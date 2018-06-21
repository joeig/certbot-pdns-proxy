package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joeig/go-powerdns"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Error struct {
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

type Config struct {
	Server struct {
		ListenAddress string `mapstructure:"listenaddress"`
		CertFile      string `mapstructure:"certFile"`
		KeyFile       string `mapstructure:"keyFile"`
	} `mapstructure:"server"`
	PowerDNS struct {
		BaseURL string `mapstructure:"baseURL"`
		VHost   string `mapstructure:"vhost"`
		APIKey  string `mapstructure:"apiKey"`
	} `mapstructure:"powerdns"`
	Miscellaneous struct {
		DefaultTTL int `mapstructure:"defaultTTL"`
	} `mapstructure:"miscellaneous"`
	AuthTable []Auth `mapstructure:"authTable"`
}

var C Config

type Auth struct {
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	FQDNRegex string `mapstructure:"fqdnRegex"`
	Domain    string `mapstructure:"domain"`
}

func (a *Auth) CheckPermission(username *string, password *string, fqdn *string) error {
	if a.Username != *username {
		return Error{Message: "Wrong username"}
	}
	if a.Password != *password {
		return Error{Message: "Wrong password"}
	}
	if matched, _ := regexp.MatchString(a.FQDNRegex, *fqdn); matched != true {
		return Error{Message: "FQDN not allowed"}
	}
	return nil
}

func checkAuthorization(r *http.Request) (string, *Auth, error) {
	fqdn := r.URL.Query().Get("fqdn")
	if fqdn == "" {
		return "", new(Auth), Error{Message: "FQDN parameter missing"}
	}

	username, password, ok := r.BasicAuth()
	if ok != true {
		return "", new(Auth), Error{Message: "Authentication header missing"}
	}

	var auth Auth
	for _, a := range C.AuthTable {
		if err := a.CheckPermission(&username, &password, &fqdn); err == nil {
			auth = a
			break
		}
	}
	if auth.Username == "" {
		return "", new(Auth), Error{Message: "No matching authentication entry"}
	}

	return fqdn, &auth, nil
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	fqdn, auth, err := checkAuthorization(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	validation := r.URL.Query().Get("validation")
	if validation == "" {
		log.Print("Validation parameter missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	validation = strings.Replace(validation, "\"", "\\\"", -1)
	validation = fmt.Sprintf("\"%s\"", validation)

	pdns := powerdns.NewClient(C.PowerDNS.BaseURL, C.PowerDNS.VHost, auth.Domain, C.PowerDNS.APIKey)
	if _, err := pdns.AddRecord(fqdn, "TXT", C.Miscellaneous.DefaultTTL, []string{validation}); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Authentication OK for %s: %s, %s", auth.Username, fqdn, validation)
	w.WriteHeader(http.StatusAccepted)
}

func Cleanup(w http.ResponseWriter, r *http.Request) {
	fqdn, auth, err := checkAuthorization(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	pdns := powerdns.NewClient(C.PowerDNS.BaseURL, C.PowerDNS.VHost, auth.Domain, C.PowerDNS.APIKey)
	if _, err := pdns.DeleteRecord(fqdn, "TXT"); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Cleanup OK for %s: %s", auth.Username, fqdn)
	w.WriteHeader(http.StatusAccepted)
}

func parseConfig(config *Config) {
	configFile := flag.String("config", "config.yaml", "Configuration file")
	flag.Parse()
	viper.SetConfigFile(*configFile)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("%s", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("%s", err))
	}
}

func runServer() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/authenticate", Authenticate).Methods("POST")
	router.HandleFunc("/v1/cleanup", Cleanup).Methods("DELETE")
	log.Fatal(http.ListenAndServeTLS(C.Server.ListenAddress, C.Server.CertFile, C.Server.KeyFile, router))
}

func main() {
	parseConfig(&C)
	runServer()
}

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
)

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

func parseConfig(config *Config, configFile *string) {
	viper.SetConfigFile(*configFile)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("%s", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("%s", err))
	}
}

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

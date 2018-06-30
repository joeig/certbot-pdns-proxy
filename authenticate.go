package main

import (
	"fmt"
	"github.com/joeig/go-powerdns"
	"log"
	"net/http"
	"strings"
	"time"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	fqdn, auth, err := checkAuthorization(r)
	if err != nil {
		log.Print("Authentication failed: ", err)
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

	cAddRecordErr := make(chan error)

	go func() {
		pdns := powerdns.NewClient(C.PowerDNS.BaseURL, C.PowerDNS.VHost, auth.Domain, C.PowerDNS.APIKey)
		_, err := pdns.AddRecord(fqdn, "TXT", C.Miscellaneous.DefaultTTL, []string{validation})
		cAddRecordErr <- err
	}()

	select {
	case c := <-cAddRecordErr:
		if c != nil {
			log.Print("Bad PowerDNS API response: ", c)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	case <-time.After(10 * time.Second):
		log.Print("PowerDNS API timeout")
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	}

	log.Printf("Authentication OK for %s: %s, %s", auth.Username, fqdn, validation)
	w.WriteHeader(http.StatusAccepted)
}

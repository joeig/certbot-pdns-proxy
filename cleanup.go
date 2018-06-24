package main

import (
	"github.com/joeig/go-powerdns"
	"log"
	"net/http"
)

func Cleanup(w http.ResponseWriter, r *http.Request) {
	fqdn, auth, err := checkAuthorization(r)
	if err != nil {
		log.Print("Authorization failed:", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	pdns := powerdns.NewClient(C.PowerDNS.BaseURL, C.PowerDNS.VHost, auth.Domain, C.PowerDNS.APIKey)
	if _, err := pdns.DeleteRecord(fqdn, "TXT"); err != nil {
		log.Print("Bad PowerDNS API response:", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	log.Printf("Cleanup OK for %s: %s", auth.Username, fqdn)
	w.WriteHeader(http.StatusAccepted)
}

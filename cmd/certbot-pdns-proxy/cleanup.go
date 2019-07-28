package main

import (
	"github.com/joeig/go-powerdns"
	"log"
	"net/http"
	"time"
)

// Cleanup handles clean up requests
func Cleanup(w http.ResponseWriter, r *http.Request) {
	fqdn, auth, err := checkAuthorization(r)
	if err != nil {
		log.Print("Authorization failed: ", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	cDeleteRecordErr := make(chan error)

	go func() {
		pdns := powerdns.NewClient(C.PowerDNS.BaseURL, C.PowerDNS.VHost, C.PowerDNS.APIKey)
		zone, _ := pdns.GetZone(auth.Domain)
		err := zone.DeleteRecord(fqdn, "TXT")
		cDeleteRecordErr <- err
	}()

	select {
	case c := <-cDeleteRecordErr:
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

	log.Printf("Cleanup OK for %s: %s", auth.Username, fqdn)
	w.WriteHeader(http.StatusAccepted)
}

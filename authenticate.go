package main

import (
	"fmt"
	"github.com/joeig/go-powerdns"
	"log"
	"net/http"
	"strings"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	fqdn, auth, err := checkAuthorization(r)
	if err != nil {
		log.Print("Authentication failed:", err)
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
		log.Print("Bad PowerDNS API response:", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	log.Printf("Authentication OK for %s: %s, %s", auth.Username, fqdn, validation)
	w.WriteHeader(http.StatusAccepted)
}

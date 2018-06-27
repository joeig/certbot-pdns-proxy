package main

import "testing"

func TestParseConfig(t *testing.T) {
	var config Config
	configFile := "config.dist.yml"
	parseConfig(&config, &configFile)
	if config.Server.ListenAddress == "" || config.PowerDNS.BaseURL == "" {
		t.Error("Config parser not working")
	}
}

func TestAuthCheckPermissionOK(t *testing.T) {
	var config Config
	configFile := "config.dist.yml"
	parseConfig(&config, &configFile)
	username, password, fqdn := "foo", "bar", "_acme-challenge.abc.example.com"
	err := config.AuthTable[0].CheckPermission(&username, &password, &fqdn)
	if err != nil {
		t.Error("Auth not OK, but should be OK")
	}
}

func TestAuthCheckPermissionNotOK(t *testing.T) {
	var config Config
	configFile := "config.dist.yml"
	parseConfig(&config, &configFile)
	username, password, fqdn := "foo", "bar", "_acme-challenge.abc.example.com"
	err := config.AuthTable[0].CheckPermission(&username, &password, &fqdn)
	if err != nil {
		t.Error("Auth OK, but it should not be OK")
	}
}

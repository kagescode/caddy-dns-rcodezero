package rcodezero

import (
	"fmt"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"

	libdnsrcodezeroacme "github.com/kagescode/libdns-rcodezeroacme"
)

func TestSingleArg(t *testing.T) {
	fmt.Println("Testing single string argument (api_token)...")

	apiToken := "abc123"
	config := fmt.Sprintf("rcodezero %s", apiToken)

	dispenser := caddyfile.NewTestDispenser(config)
	p := Provider{&libdnsrcodezeroacme.Provider{}}

	err := p.UnmarshalCaddyfile(dispenser)
	if err != nil {
		t.Fatalf("UnmarshalCaddyfile failed with %v", err)
	}

	if p.Provider.APIToken != apiToken {
		t.Fatalf("Expected APIToken to be %q but got %q", apiToken, p.Provider.APIToken)
	}
}

func TestAPITokenInBlock(t *testing.T) {
	fmt.Println("Testing api_token provided in block...")

	apiToken := "abc123"
	config := fmt.Sprintf(`rcodezero {
		api_token %s
	}`, apiToken)

	dispenser := caddyfile.NewTestDispenser(config)
	p := Provider{&libdnsrcodezeroacme.Provider{}}

	err := p.UnmarshalCaddyfile(dispenser)
	if err != nil {
		t.Fatalf("UnmarshalCaddyfile failed with %v", err)
	}

	if p.Provider.APIToken != apiToken {
		t.Fatalf("Expected APIToken to be %q but got %q", apiToken, p.Provider.APIToken)
	}
}

func TestBaseURLInBlock(t *testing.T) {
	fmt.Println("Testing base_url provided in block...")

	apiToken := "abc123"
	baseURL := "https://my.rcodezero.at"
	config := fmt.Sprintf(`rcodezero {
		api_token %s
		base_url  %s
	}`, apiToken, baseURL)

	dispenser := caddyfile.NewTestDispenser(config)
	p := Provider{&libdnsrcodezeroacme.Provider{}}

	err := p.UnmarshalCaddyfile(dispenser)
	if err != nil {
		t.Fatalf("UnmarshalCaddyfile failed with %v", err)
	}

	if p.Provider.APIToken != apiToken {
		t.Fatalf("Expected APIToken to be %q but got %q", apiToken, p.Provider.APIToken)
	}
	if p.Provider.BaseURL != baseURL {
		t.Fatalf("Expected BaseURL to be %q but got %q", baseURL, p.Provider.BaseURL)
	}
}

func TestEmptyConfig(t *testing.T) {
	fmt.Println("Testing empty config fails to parse...")

	config := "rcodezero"

	dispenser := caddyfile.NewTestDispenser(config)
	p := Provider{&libdnsrcodezeroacme.Provider{}}

	err := p.UnmarshalCaddyfile(dispenser)
	if err == nil {
		t.Fatalf("Expected error, got none (api_token=%q base_url=%q)", p.Provider.APIToken, p.Provider.BaseURL)
	}
}

func TestTooManyArgs(t *testing.T) {
	fmt.Println("Testing too many args...")

	apiToken := "abc123"
	config := fmt.Sprintf("rcodezero %s extra", apiToken)

	dispenser := caddyfile.NewTestDispenser(config)
	p := Provider{&libdnsrcodezeroacme.Provider{}}

	err := p.UnmarshalCaddyfile(dispenser)
	if err == nil {
		t.Fatalf("Expected error for too many args, got none")
	}
}

func TestDuplicateToken(t *testing.T) {
	fmt.Println("Testing duplicate api_token fails...")

	config := `rcodezero {
		api_token abc
		api_token def
	}`

	dispenser := caddyfile.NewTestDispenser(config)
	p := Provider{&libdnsrcodezeroacme.Provider{}}

	err := p.UnmarshalCaddyfile(dispenser)
	if err == nil {
		t.Fatalf("Expected error for duplicate api_token, got none")
	}
}

func TestProvisionMissingTokenFails(t *testing.T) {
	fmt.Println("Testing Provision() fails with missing token...")

	p := Provider{&libdnsrcodezeroacme.Provider{}}
	err := p.Provision(caddy.Context{})
	if err == nil {
		t.Fatalf("Expected Provision to fail with missing token, but it succeeded")
	}
}

func TestProvisionTokenOK(t *testing.T) {
	fmt.Println("Testing Provision() succeeds with token...")

	p := Provider{&libdnsrcodezeroacme.Provider{
		APIToken: "abc123",
	}}
	err := p.Provision(caddy.Context{})
	if err != nil {
		t.Fatalf("Expected Provision to succeed, got error: %v", err)
	}
}


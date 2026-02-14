package rcodezero

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"

	libdnsrcodezeroacme "github.com/kagescode/libdns-rcodezeroacme"
)

// Provider lets Caddy solve the ACME DNS challenge by manipulating DNS records
// via the RcodeZero ACME endpoint (through libdns-rcodezeroacme).
type Provider struct {
	*libdnsrcodezeroacme.Provider
}

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "dns.providers.rcodezero",
		New: func() caddy.Module {
			return &Provider{new(libdnsrcodezeroacme.Provider)}
		},
	}
}

// Provision allows placeholder replacement (e.g. {$ENV_VAR}) in config values.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()

	p.Provider.APIToken = repl.ReplaceAll(p.Provider.APIToken, "")
	p.Provider.BaseURL = repl.ReplaceAll(p.Provider.BaseURL, "")

	// Validate required config
	if p.Provider.APIToken == "" {
		return fmt.Errorf("rcodezero: missing api token")
	}

	// BaseURL is optional; libdns-rcodezeroacme can default internally.
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens.
//
// Syntax:
//
//	rcodezero {
//	  api_token <token>
//	  base_url  <url>     # optional
//	}
//
// Also supported:
//
//	rcodezero <token>
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// Optional inline token: rcodezero <token>
		if d.NextArg() {
			p.Provider.APIToken = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if p.Provider.APIToken != "" {
					return d.Err("API token already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.APIToken = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}

			case "base_url":
				if p.Provider.BaseURL != "" {
					return d.Err("Base URL already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.BaseURL = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}

			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}

	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}

	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)


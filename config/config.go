package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

// FromCLIOpts builds a Config from command line options and env vars.
// The urfave/cli Context gives us access to both on program start, and also
// allows us to set defaults.
func FromCLIOpts(ctx *cli.Context) Config {
	var c Config
	c.DBPath = ctx.String("db")
	c.APICert = ctx.String("apiCert")
	c.APIKey = ctx.String("apiKey")
	c.WebUICert = ctx.String("webUICert")
	c.WebUIKey = ctx.String("webUIKey")
	c.ProxyCert = ctx.String("proxyCert")
	c.ProxyKey = ctx.String("proxyKey")
	c.Auth0ClientID = ctx.String("auth0ClientID")
	c.Auth0Secret = ctx.String("auth0Secret")
	c.BypassAuth0 = ctx.BoolT("bypassAuth0")
	return c
}

// The Config struct. Contains config values suitable for multiple processes we
// spin off, including the pure gRPC management API, the http(s) listener for
// our GopherJS Web UI (including websockets), and the listener for our proxy
// forwarding implementation.
type Config struct {
	// Filepath to the boltdb file.
	DBPath string `toml:"db_path"`

	// APICert and APIKey are paths to PEM-encoded TLS assets
	// for our pure gRPC api
	APICert string `toml:"api_cert"`
	APIKey  string `toml:"api_key"`

	// WebUICert and WebUIKey are paths to PEM-encoded TLS assets
	// for our GopherJS-over-websockets UI. This UI is wrapped with
	// auth0 handlers.
	WebUICert string `toml:"webui_cert"`
	WebUIKey  string `toml:"webui_key"`

	// ProxyCert and  are paths to PEM-encoded TLS assets
	// for our GopherJS-over-websockets UI. This UI is wrapped with
	// auth0 handlers.
	ProxyCert string `toml:"proxy_cert"`
	ProxyKey  string `toml:"proxy_key"`

	// Auth0 config values
	Auth0ClientID string `toml:"auth0_client_id"`
	Auth0Secret   string `toml:"auth0_secret"`
	BypassAuth0   bool   `toml:"bypass_auth0"`
}

// ExampleConfig can be written to disk. See the systemd-install command.
var ExampleConfig = `

	# Filepath to the boltdb file.
	db_path = ""

	# Paths to PEM-encoded TLS assets
	# for our pure gRPC api
	api_cert = ""
	api_key = ""

	# WebUICert and WebUIKey are paths to PEM-encoded TLS assets
	# for our GopherJS-over-websockets UI. This UI is wrapped with
	# auth0 handlers.
	webui_cert = ""
	webui_key = ""

	# ProxyCert and  are paths to PEM-encoded TLS assets
	# for our GopherJS-over-websockets UI. This UI is wrapped with
	# auth0 handlers.
	proxy_cert = ""
	proxy_key = ""

	# Auth0 config values
	auth0_client_id = ""
	auth0_secret = ""
	bypass_auth0 = false

`

// UnitFile can be written to /etc/systemd/system/co-chair.service
// to aid in installing co-chair as a linux service.
var UnitFile = `
[Unit]
Description=co-chair configurable proxy
After=network.target

[Service]
ExecStart=/usr/local/bin/co-chair serve --conf /opt/co-chair/conf.toml 
RestartSec=3
Restart=on-failure

# Every 10 min, try to restart the dead service.
StartLimitInterval=10min

# We can only fail 5 times within the 10 min interval.
StartLimitBurst=5

StartLimitAction=none

[Install]
WantedBy=multi-user.target

`

// SystemDInstall places a config file at /opt/co-chair/conf.toml and a systemd
// unit file at /etc/systemd/system/co-chair.service.
func SystemDInstall(conf Config) error {

	unitFile := "/etc/systemd/system/co-chair.service"
	if exists, err := exists(unitFile); exists || err != nil {
		if err != nil {
			return fmt.Errorf("os stat: %v", err)
		}
		return fmt.Errorf("file exists: %s", unitFile)
	}
	if err := ioutil.WriteFile(unitFile, []byte(UnitFile), 0644); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	if err := os.MkdirAll("/opt/co-chair", 0644); err != nil {
		return fmt.Errorf("mkdir all: %v", err)
	}

	exampleConf := "/opt/co-chair/conf.toml"
	if exists, err := exists(exampleConf); exists || err != nil {
		if err != nil {
			return fmt.Errorf("os stat: %v", err)
		}
		return fmt.Errorf("file exists: %s", exampleConf)
	}
	if err := ioutil.WriteFile(exampleConf, []byte(ExampleConfig), 0644); err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

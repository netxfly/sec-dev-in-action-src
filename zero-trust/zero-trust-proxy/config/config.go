package config

import (
	"time"
)

// Config structure used to configure Helios
type Config struct {
	Server    Server     `yaml:"server"`
	Upstreams []Upstream `yaml:"upstreams"`
	Routes    []Route    `yaml:"routes"`
	Identity  Identity   `yaml:"identity"`
	JWT       JWT        `yaml:"jwt"`
}

// Server structure is used to configure the HTTP(S) server
type Server struct {
	ListenIP    string        `yaml:"listen_ip"`
	ListenPort  int           `yaml:"listen_port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	TLSContext  TLSContext    `yaml:"tls_context"`
}

// TLSContext structure is used to configure TLS for the server
type TLSContext struct {
	CertificatePath string `yaml:"certificate_path"`
	PrivateKeyPath  string `yaml:"private_key_path"`
}

// Route represents a route configuration
type Route struct {
	Host  string
	Rules []string
	HTTP  struct {
		Paths []struct {
			Path           string
			Upstream       string
			Authentication bool `yaml:"authentication"`
		}
	}
}

// Upstream represents a single proxy upstream
type Upstream struct {
	Name           string `yaml:"name"`
	URL            string `yaml:"url"`
	ConnectTimeout time.Duration
}

// Identity provider configuration
type Identity struct {
	Provider     string `yaml:"provider"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	OAuth2       struct {
		AuthURL    string `yaml:"auth_url"`
		TokenURL   string `yaml:"token_url"`
		ProfileURL string `yaml:"profile_url"`
	}
}

// JWT token configuration
type JWT struct {
	Secret  string
	Expires time.Duration
}

// UnmarshalYAML parses upstream configuration from a YAML file
func (c *Upstream) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	buf := struct {
		ConnectTimeout string `yaml:"connect_timeout"`
		Name           string `yaml:"name"`
		URL            string `yaml:"url"`
	}{}

	if err := unmarshal(&buf); err != nil {
		return err
	}

	timeout, err := time.ParseDuration(buf.ConnectTimeout)
	if err != nil {
		return err
	}

	c.ConnectTimeout = timeout
	c.URL = buf.URL
	c.Name = buf.Name

	return nil
}

// UnmarshalYAML parses JWT configuration from a YAML file
func (c *JWT) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	buf := struct {
		Secret  string `yaml:"secret"`
		Expires string `yaml:"expires"`
	}{}

	if err := unmarshal(&buf); err != nil {
		return err
	}

	expires, err := time.ParseDuration(buf.Expires)
	if err != nil {
		return err
	}

	c.Expires = expires
	c.Secret = buf.Secret

	return nil
}

// UnmarshalYAML parses server configuration from a YAML file
func (c *Server) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var buf struct {
		ListenIP    string     `yaml:"listen_ip"`
		ListenPort  int        `yaml:"listen_port"`
		Timeout     string     `yaml:"timeout"`
		IdleTimeout string     `yaml:"idle_timeout"`
		TLSContext  TLSContext `yaml:"tls_context"`
	}

	if err := unmarshal(&buf); err != nil {
		return err
	}

	timeout, err := time.ParseDuration(buf.Timeout)
	if err != nil {
		return err
	}

	idleTimeout, err := time.ParseDuration(buf.IdleTimeout)
	if err != nil {
		return err
	}

	c.Timeout = timeout
	c.IdleTimeout = idleTimeout
	c.TLSContext = buf.TLSContext
	c.ListenIP = buf.ListenIP
	c.ListenPort = buf.ListenPort

	return nil
}

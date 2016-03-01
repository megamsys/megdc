package subd

import (
  "errors"
  "os"
	"os/user"
	"path/filepath"

  "github.com/BurntSushi/toml"
  log "github.com/Sirupsen/logrus"
)

const (
	NAME     = "vertice-prog"

  DefaultPath = "/var/lib/megam"

)

type Settings struct {
	Name 			   string
	Hosts        `toml:"common"`
  One    			 `toml:one`
	Dns  				 `toml:dns`
	Docker 			 `toml:docker`
  Organise     `toml:org`
}


func (c Settings) String() string {
	return ("\n" + c.Hosts.String() + "\n" +
"\n" + c.One.String() + "\n" +
"\n" + c.Dns.String() + "\n" +
"\n" + c.Docker.String() + "\n" +
"\n" + c.Organise.String() + "\n")
}

// NewConfig returns an instance of Config with reasonable defaults.
func NewConfig() *Settings {
	c := &Settings{
    Name: NAME,
  }
	c.Hosts = HostConfig()
  c.Dns = DnsConfig()
  c.Docker = DockerConfig()
  c.One = OneConfig()
  c.Organise = OrgConfig()
	return c
}

// Validate returns an error if the config is invalid.
func (c *Settings) Validate() error {
	if c.Hosts.Home == "" {
		return errors.New("Home Dir must be specified")
	}
	return nil
}

func ParseConfig() (*Settings, error) {

  var path string
	if os.Getenv("MEGAM_HOME") != "" {
		path = os.Getenv("MEGAM_HOME")
	} else if u, err := user.Current(); err == nil {
		path = u.HomeDir
	}

	path = filepath.Join(path, "/megdc/megdc.conf")
	config := NewConfig()
	if path == "" {
		path = DefaultPath + "/megdc/megdc.conf"
	}
	log.Warnf("Using configuration at: %s", path)
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	log.Debug(config)
	return config, nil
}

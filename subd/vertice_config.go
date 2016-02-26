package subd

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"github.com/megamsys/libgo/cmd"
)

const (
	// DefaultRiak is the default riak if one is not provided.
	DefaultScylla = "127.0.0.1"

	// DefaultApi is the default megam gateway if one is not provided.
	DefaultApi = "http://localhost:9000/v2/"

	// DefaultNSQ is the default nsqd if its not provided.
	DefaultNSQd = "localhost:4151"
  DefaultName = "vertice-dev"
	//default user
	DefaultUser = "megam"

	MEGAM_HOME = "MEGAM_HOME"
)


// Config represents the meta configuration.
type Config struct {
	Name			 string 								 `toml:"name"`
	Home       string                  `toml:"home"`
	Scylla     string         				 `toml:"scylla"`
	NSQd       string                  `toml:"nsqd"`
	Api        string                  `toml:"api"`
}

var MC *Config

func (c Config) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Config:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Meta", "cyan", "", "") + "\n"))
	b.Write([]byte("Name      " + "\t" + c.Name + "\n"))
	b.Write([]byte("Home      " + "\t" + c.Home + "\n"))
	b.Write([]byte("Scylla      " + "\t" + c.Scylla + "\n"))
	b.Write([]byte("Api       " + "\t" + c.Api + "\n"))
	b.Write([]byte("NSQd      " + "\t" + c.NSQd + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func NewConfig() *Config {
	var homeDir string
	// By default, store logs, meta and load conf files in MEGAM_HOME directory
	if os.Getenv(MEGAM_HOME) != "" {
		homeDir = os.Getenv(MEGAM_HOME)
	} else if u, err := user.Current(); err == nil {
		homeDir = u.HomeDir
	} else {
		return nil
	}

	defaultDir := filepath.Join(homeDir, "megdc/")

	// Config represents the configuration format for the vertice.
	return &Config{
		Name: DefaultName,
		Home: defaultDir,
		Scylla: DefaultScylla,
		Api:  DefaultApi,
		NSQd: DefaultNSQd,
	}
}

func (c *Config) MkGlobal() {
	MC = c
}

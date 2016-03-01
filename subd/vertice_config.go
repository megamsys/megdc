package subd

import (
	"bytes"
	"fmt"

	"strings"
	"text/tabwriter"
	"github.com/megamsys/libgo/cmd"
)

const (
	// DefaultHost
	localhost = "127.0.0.1"
	// DefaultRiak is the default Scylla if one is not provided.

	DefaultScylla = localhost

	// DefaultApi is the default megam gateway if one is not provided.
	DefaultApi = "http://localhost:9000/v2/"

	// DefaultNSQ is the default nsqd if its not provided.
	DefaultNSQd = localhost

  DefaultName = "vertice-dev"
	//default user
	DefaultUser = "megam"

  DefaultHome = "/var/lib/megam/"

	MEGAM_HOME = "MEGAM_HOME"

	DefaultSwarmMaster = localhost

	DefaultAccessKey = "dummy-access-key"

	DefaultScretKey  = "dummy-scret-key"

	DefaultEndpoint = localhost

	DefaultOneUserId = "oneadmin"

	DefaultOnePassword = "dummy-password"

	DefaultOrg = "megam.io"

	DefaultDomain = "example.com"

)


// Config represents the meta configuration.
type Hosts struct {
	Home       string                  `toml:"home"`
	Scylla     string         				 `toml:"scylla"`
	NSQd       string                  `toml:"nsqd"`
	Api        string                  `toml:"api"`
}


type Dns struct{
  Route53_access_key string `toml:route53_access`
  Route53_secret_key string `toml:route53_secret`
}

type Docker struct{
  Docker_swarm  string     `toml:swarm`
}

type One struct{
  One_endpoint  string  `toml:one_endpoint`
  One_userid    string  `toml:one_userid`
  One_password  string  `toml:one_password`
}

type Organise struct{
	Org 		string `toml:org`
	Domain  string `toml:domain`
}

func (c Hosts) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Hosts:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Meta", "cyan", "", "") + "\n"))
	b.Write([]byte("Home      " + "\t" + c.Home + "\n"))
	b.Write([]byte("Scylla      " + "\t" + c.Scylla + "\n"))
	b.Write([]byte("Api       " + "\t" + c.Api + "\n"))
	b.Write([]byte("NSQd      " + "\t" + c.NSQd + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c Dns) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("DNS:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Meta", "cyan", "", "") + "\n"))
	b.Write([]byte("Route53 AccessKey      " + "\t" + c.Route53_access_key + "\n"))
	b.Write([]byte("Route53 SecretKey      " + "\t" + c.Route53_secret_key + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c One) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("One:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Meta", "cyan", "", "") + "\n"))
	b.Write([]byte("One end point     " + "\t" + c.One_endpoint + "\n"))
	b.Write([]byte("One User Id      " + "\t" + c.One_userid + "\n"))
	b.Write([]byte("One Password      " + "\t" + c.One_password + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c Docker) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Docker:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Meta", "cyan", "", "") + "\n"))
	b.Write([]byte("Docker Swarm Master      " + "\t" + c.Docker_swarm + "\n"))
	b.Write([]byte("---\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func (c Organise) String() string {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Organisation:", "white", "", "bold") + "\t" +
		cmd.Colorfy("Meta", "cyan", "", "") + "\n"))
	b.Write([]byte("Name      " + "\t" + c.Org + "\n"))
	b.Write([]byte("Domain      " + "\t" + c.Domain + "\n"))
	b.Write([]byte("---\n\n"))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String())
}

func DnsConfig() Dns{
	return Dns{
		Route53_access_key: DefaultAccessKey,
		Route53_secret_key: DefaultScretKey,
	}
}

func OneConfig() One{
	return One{
    One_endpoint: DefaultEndpoint,
    One_userid: DefaultOneUserId,
    One_password: DefaultOnePassword,
	}
}

func DockerConfig() Docker{
	return Docker{
		Docker_swarm: DefaultSwarmMaster,
	}
}
func HostConfig() Hosts {
		// Config represents the configuration format for the vertice.
	return Hosts{
		Home: DefaultHome,
		Scylla: DefaultScylla,
		Api:  DefaultApi,
		NSQd: DefaultNSQd,
	}
}

func OrgConfig() Organise {
	return Organise{
		Org: DefaultOrg,
		Domain: DefaultDomain,
	}
}

package subd

import (
//  "fmt"
  "testing"

	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{

}

var _ = check.Suite(&S{})

// Ensure the configuration can be parsed.
func (s *S) TestMetaConfig_Parse(c *check.C) {
	// Parse configuration.
	var cm Hosts
	if _, err := toml.Decode(`
home = "/var/lib/megam/megdc"
api = "https://api.megam.io"
nsqd = "localhost:4150"
`, &cm); err != nil {
		c.Fatal(err)
	}
	c.Assert(cm.Home, check.Equals, "/var/lib/megam/megdc")
	c.Assert(cm.Api, check.Equals, "https://api.megam.io")
	c.Assert(cm.NSQd, check.DeepEquals, "localhost:4150")
}

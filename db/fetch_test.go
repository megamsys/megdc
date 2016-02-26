package db

import (
  "fmt"
  "testing"

	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}
type Ss struct{}

type Sub struct{
  Name       string                  `toml:"name"`
  Home       string                  `toml:"home"`
	Scylla     string         				 `toml:"scylla"`
	NSQd       string                  `toml:"nsqd"`
	Api        string                  `toml:"api"`
}

func testConfig() *Sub {
  return &Sub{
		Home: "1",
		Scylla: "2",
		Api:  "3",
		NSQd: "4",
	}
}


type Mconf struct{
  Common *Sub  `toml:"common"`
}

func testConfigToml() *Mconf {
	c := &Mconf{}
	c.Common = testConfig()
	return c
}

var _ = check.Suite(&Ss{})

// Ensure the configuration can be parsed.
func (s *S) TestMetaConfig_Parse(c *check.C)  {
	// Parse configuration.
	cm := testConfigToml()
  path :="/home/megam/code/workspace/yaml/megdc/megdc.conf"
	if _, err := toml.DecodeFile(path, &cm); err != nil {
		fmt.Println(err)
	}
	c.Assert(cm.Common.Home, check.Equals, "/var/lib/megam/s")
}

/*func (s *S) TestStore(c *check.C) {
fmt.Println("*****************")
 fmt.Println(cm.Common)
 st := &Sub{}
 t := TableInfo{
   Name: "testing",
   Pks: []string{"name"},
   Ccms: []string{},
   Query: map[string]string{"name": "vertice-dev"},
 }
 err := Write(t, st)
 fmt.Println("*******************************")
 fmt.Println(err)
//  c.Assert(nil, check.NotNil)
}*/

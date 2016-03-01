package db

/*
import (
  "fmt"
  "testing"

	"github.com/BurntSushi/toml"
	"gopkg.in/check.v1"
)

const (
  Path ="/home/megam/code/workspace/yaml/megdc/megdc.conf"
)

func Test(t *testing.T) {
	check.TestingT(t)
}
type S struct{
  sc  Mconf
}

type Tcommon struct{
  Home       string                  `toml:"home"`
	Scylla     string         				 `toml:"scylla"`
	NSQd       string                  `toml:"nsqd"`
	Api        string                  `toml:"api"`
}

func testConfigcommon() Tcommon {
  return Tcommon{
		Home: "1",
		Scylla: "2",
		Api:  "3",
		NSQd: "4",
	}
}

func testConfigOne() TOne {
  return TOne{
    One_endpoint: "127.0.0.1",
    One_id: "oneadmin",
    One_password: "onesdfpassword",
	}
}

type Mconf struct{
  Name       string     `toml:"name"`
  Tcommon  `toml:"common"`
  TOne        `toml:"one"`
}

type dns struct{
  Route53_access string `toml:"route53_access"`
  Route53_secret string `toml:"route53_secret"`
}

type docker struct{
 Swarm  string     `toml:swarm`
}

type TOne struct{
  One_endpoint  string  `toml:end_point`
  One_id        string  `toml:one_userid`
  One_password  string  `toml:one_password`
}

func testConfigToml() *Mconf {
	c := &Mconf{
    Name: "dev-prog",
  }
	c.Tcommon = testConfigcommon()
  c.TOne = testConfigOne()
	return c
}

var _ = check.Suite(&S{})

// Ensure the configuration can be parsed.
func (s *S) TestMetaConfig_ParseCommon(c *check.C)  {
	// Parse configuration.
  cm := testConfigcommon()
	if _, err := toml.DecodeFile(Path, &cm); err != nil {
		fmt.Println(err)
	}

  t := TableInfo{
    Name: "test_common",
    Pks: []string{"Home"},
    Ccms: []string{},
  }
  err := Write(t, cm)
  c.Assert(err,check.IsNil)
  fmt.Printf("%#v",cm)
//	c.Assert(cm.Home, check.Equals,"/var/lib/megam/")
}

func (s *S) TestMetaConfig_Parse(c *check.C)  {
	// Parse configuration.
	cm := testConfigToml()
	if _, err := toml.DecodeFile(Path, &cm); err != nil {
		fmt.Println(err)
	}
  t := TableInfo{
    Name: "test_db",
    Pks: []string{"Name"},
    Ccms: []string{},
  }

  fmt.Println(cm,"\n\n")

  err := Write(t, cm)
	c.Assert(err, check.IsNil)
}
*/

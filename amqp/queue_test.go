package amqp

import (
	"github.com/tsuru/config"
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

type S struct{}

var _ = gocheck.Suite(&S{})

func (s *S) TestMessageDelete(c *gocheck.C) {
	m := Message{}
	c.Assert(m.delete, gocheck.Equals, false)
	m.Delete()
	c.Assert(m.delete, gocheck.Equals, true)
}

func (s *S) TestFactory(c *gocheck.C) {
	config.Set("queue", "rabbitmq")
	defer config.Unset("queue")
	f, err := Factory()
	c.Assert(err, gocheck.IsNil)
	_, ok := f.(rabbitmqFactory)
	c.Assert(ok, gocheck.Equals, true)
}

func (s *S) TestFactoryConfigUndefined(c *gocheck.C) {
	f, err := Factory()
	c.Assert(err, gocheck.IsNil)
	_, ok := f.(rabbitmqFactory)
	c.Assert(ok, gocheck.Equals, true)
}

func (s *S) TestFactoryConfigUnknown(c *gocheck.C) {
	config.Set("queue", "unknown")
	defer config.Unset("queue")
	f, err := Factory()
	c.Assert(f, gocheck.IsNil)
	c.Assert(err, gocheck.NotNil)
	c.Assert(err.Error(), gocheck.Equals, `Queue "unknown" is not known.`)
}

func (s *S) TestRegister(c *gocheck.C) {
	config.Set("queue", "unregistered")
	defer config.Unset("queue")
	Register("unregistered", rabbitmqFactory{})
	_, err := Factory()
	c.Assert(err, gocheck.IsNil)
}

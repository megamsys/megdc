package main

import (
	"github.com/megamsys/cloudinabox/cmd"
	"launchpad.net/gocheck"
)

func (s *S) TestCommandsFromBaseManagerAreRegistered(c *gocheck.C) {
	baseManager := cmd.BuildBaseManager("megam", version, header)
	manager := buildManager("megam")
	for name, instance := range baseManager.Commands {
		command, ok := manager.Commands[name]
		c.Assert(ok, gocheck.Equals, true)
		c.Assert(command, gocheck.FitsTypeOf, instance)
	}
}

func (s *S) TestAppStartIsRegistered(c *gocheck.C) {
	manager := buildManager("megam")
	create, ok := manager.Commands["startapp"]
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(create, gocheck.FitsTypeOf, &AppStart{})
}

func (s *S) TestAppStopIsRegistered(c *gocheck.C) {
	manager := buildManager("megam")
	remove, ok := manager.Commands["stopapp"]
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(remove, gocheck.FitsTypeOf, &AppStop{})
}
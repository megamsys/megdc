package main

import (
	"github.com/megamsys/libgo/cmd"
	"launchpad.net/gocheck"
)

func (s *S) TestCIBNodeStart(c *gocheck.C) {
	desc := `starts the cib node daemon.


`
	expected := &cmd.Info{
		Name:    "startnode",
		Usage:   `startnode`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := CIBNodeStart{}
	c.Assert(command.Info(), gocheck.DeepEquals, expected)
}

package main

import (
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
)

func (s *S) TestCIBNodeStart(c *check.C) {
	desc := `starts the cib node daemon.


`
	expected := &cmd.Info{
		Name:    "startnode",
		Usage:   `startnode`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := CIBNodeStart{}
	c.Assert(command.Info(), check.DeepEquals, expected)
}

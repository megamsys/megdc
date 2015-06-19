package main

import (
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
)

func (s *S) TestCIBStartInfo(c *check.C) {
	desc := `starts the cib base web daemon.
	
	
	`
	expected := &cmd.Info{
		Name:    "start",
		Usage:   `start`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := CIBStart{}
	c.Assert(command.Info(), check.DeepEquals, expected)
}
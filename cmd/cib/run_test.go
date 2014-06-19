package main

import (
	"github.com/megamsys/cloudinabox/cmd"
	"launchpad.net/gocheck"
)

func (s *S) TestGulpStartInfo(c *gocheck.C) {
	desc := `starts the gulpd daemon, and connects to queue.

If you use the '--dry' flag gulpd will do a dry run(parse conf/jsons) and exit.

`

	expected := &cmd.Info{
		Name:    "start",
		Usage:   `start [--dry] [--config]`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := GulpStart{}
	c.Assert(command.Info(), gocheck.DeepEquals, expected)
}


func (s *S) TestGulpStopInfo(c *gocheck.C) {
	desc := `stops the gulpd daemon, and shutsdown the queue.

If you use the '--bark' flag gulpd will notify daemon status.

`
	expected := &cmd.Info{
		Name:    "stop",
		Usage:   `stop [--bark]`,
		Desc:    desc,
		MinArgs: 0,
	}
	command := GulpStop{}
	c.Assert(command.Info(), gocheck.DeepEquals, expected)
}


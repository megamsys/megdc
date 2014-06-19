package main

import (
	"bytes"
	"github.com/megamsys/cloudinabox/cmd"
	"github.com/megamsys/cloudinabox/cmd/testing"
//	"io/ioutil"
	"launchpad.net/gocheck"
	"net/http"
//	"strings"
)

func (s *S) TestAppStartInfo(c *gocheck.C) {
	expected := &cmd.Info{
		Name:    "startapp",
		Usage:   "startapp <appname> <lifecycle_when>",
		Desc:    "starts the installed app.",
		MinArgs: 1,
	}
	c.Assert((&AppStart{}).Info(), gocheck.DeepEquals, expected)
}

func (s *S) TestAppStart(c *gocheck.C) {
	var stdout, stderr bytes.Buffer
//	result := `{"status":"success", "repository_url":"git@github.com/megamsys:nilavu.git"}`
	expected := `App "ble.megam.co" is being started!
Use appreqs list to check the status of the app.` + "\n"
	context := cmd.Context{
		Args:   []string{"ble.megam.co", "rails"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
/*	trans := testing.ConditionalTransport{
		Transport: testing.Transport{Message: result, Status: http.StatusOK},
		CondFunc: func(req *http.Request) bool {
			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			c.Assert(err, gocheck.IsNil)
			c.Assert(string(body), gocheck.Equals, `{"name":"ble","platform":"django"}`)
			return req.Method == "POST" && req.URL.Path == "/apps"
		},
	}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	*/
	command := AppStart{}
	trans := testing.Transport{Message: "success", Status: http.StatusOK}
	client := cmd.NewClient(&http.Client{Transport: &trans}, nil, manager)
	err := command.Run(&context, client)
	c.Assert(err, gocheck.IsNil)
	c.Assert(stdout.String(), gocheck.Equals, expected)
}

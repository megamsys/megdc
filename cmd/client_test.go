package cmd

import (
	"bytes"
	ttesting "github.com/indykish/gulp/cmd/testing"
	"launchpad.net/gocheck"
	"net/http"
)

func (s *S) TestShouldSetCloseToTrue(c *gocheck.C) {
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, gocheck.IsNil)
	transport := ttesting.Transport{
		Status:  http.StatusOK,
		Message: "OK",
	}
	client := NewClient(&http.Client{Transport: &transport}, nil, manager)
	client.Do(request)
	c.Assert(request.Close, gocheck.Equals, true)
}

func (s *S) TestShouldReturnBodyMessageOnError(c *gocheck.C) {
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, gocheck.IsNil)
	client := NewClient(&http.Client{Transport: &ttesting.Transport{Message: "You must be authenticated to execute this command.", Status: http.StatusUnauthorized}}, nil, manager)
	response, err := client.Do(request)
	c.Assert(response, gocheck.NotNil)
	c.Assert(err, gocheck.NotNil)
	c.Assert(err.Error(), gocheck.Equals, "You must be authenticated to execute this command.")
}



func (s *S) TestShouldSkipValidationIfThereIsNoSupportedHeaderDeclared(c *gocheck.C) {
	var buf bytes.Buffer
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, gocheck.IsNil)
	context := Context{
		Stderr: &buf,
	}
	trans := ttesting.Transport{Message: "", Status: http.StatusOK, Headers: map[string][]string{"Supported-Tsuru": {"0.3"}}}
	manager := Manager{
		version: "0.2.1",
	}
	client := NewClient(&http.Client{Transport: &trans}, &context, &manager)
	_, err = client.Do(request)
	c.Assert(err, gocheck.IsNil)
	c.Assert(buf.String(), gocheck.Equals, "")
}

func (s *S) TestShouldSkupValidationIfServerDoesNotReturnSupportedHeader(c *gocheck.C) {
	var buf bytes.Buffer
	request, err := http.NewRequest("GET", "/", nil)
	c.Assert(err, gocheck.IsNil)
	context := Context{
		Stderr: &buf,
	}
	trans := ttesting.Transport{Message: "", Status: http.StatusOK}
	manager := Manager{
		name:          "glb",
		version:       "0.2.1",
		versionHeader: "Supported-Tsuru",
	}
	client := NewClient(&http.Client{Transport: &trans}, &context, &manager)
	_, err = client.Do(request)
	c.Assert(err, gocheck.IsNil)
	c.Assert(buf.String(), gocheck.Equals, "")
}

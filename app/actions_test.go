package app

import (
/*	"errors"
	"fmt"
	"github.com/tsuru/config"
	"github.com/indykish/gulp/action"
	"launchpad.net/gocheck"
	"sort"
	"strings"*/
)

/*
func (s *S) TestCreateRepositoryForwardInvalidType(c *gocheck.C) {
	ctx := action.FWContext{Params: []interface{}{"something"}}
	_, err := createRepository.Forward(ctx)
	c.Assert(err, gocheck.NotNil)
	c.Assert(err.Error(), gocheck.Equals, "First parameter must be App or *App.")
}

func (s *S) TestCreateRepositoryBackward(c *gocheck.C) {
	h := testHandler{}
	ts := s.t.StartGandalfTestServer(&h)
	defer ts.Close()
	app := App{Name: "someapp"}
	ctx := action.BWContext{FWResult: &app, Params: []interface{}{app}}
	createRepository.Backward(ctx)
	c.Assert(h.url[0], gocheck.Equals, "/repository/someapp")
	c.Assert(h.method[0], gocheck.Equals, "DELETE")
	c.Assert(string(h.body[0]), gocheck.Equals, "null")
}
*/

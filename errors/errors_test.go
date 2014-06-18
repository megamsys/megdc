package errors

import (
	"launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct{}

var _ = gocheck.Suite(&S{})


func (s *S) TestValidationError(c *gocheck.C) {
	e := ValidationError{Message: "something"}
	c.Assert(e.Error(), gocheck.Equals, "something")
}

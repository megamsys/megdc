/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
*/
package main

import (
	"github.com/megamsys/libgo/cmd"
	"gopkg.in/check.v1"
	"github.com/megamsys/megdc/packages/megam"
//	"github.com/megamsys/megdc/packages/ceph"
//	"github.com/megamsys/megdc/packages/one"
	"os"
)


type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestCommandsFromBaseManagerAreRegistered(c *check.C) {
	baseManager := cmd.NewManager("megdc", "0.9.1", os.Stdout, os.Stderr, os.Stdin, nil, nil)
	manager := cmdRegistry("megdc")

	for name, instance := range baseManager.Commands {
		command, ok := manager.Commands[name]
		c.Assert(ok, check.Equals, true)
		c.Assert(command, check.FitsTypeOf, instance)
	}

}

func (s *S) TestStartIsRegistered(c *check.C) {
	manager := cmdRegistry("megdc")
	create, ok := manager.Commands["megaminstall"]
	c.Assert(ok, check.Equals, true)
	c.Assert(create, check.FitsTypeOf, &megam.Megaminstall{})
}



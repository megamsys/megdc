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
package megam

import (
	"fmt"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
	"launchpad.net/gnuflag"
	"reflect"
	//	"strconv"
)

type Megaminstall struct {
	Fs           *gnuflag.FlagSet
	All          bool
	Megamnilavu  bool
	Megamgateway bool
	Megamd       bool
	Megamcommon  bool
	Megammonitor bool
	Host		 string
	Username	 string
	Password     string
	Quiet        bool
}

func (g *Megaminstall) Info() *cmd.Info {
	desc := `starts megdc.

If you use the '--quiet' flag megdc doesn't print the logs.

`
	return &cmd.Info{
		Name:    "megaminstall",
		Usage:   `megaminstall [--all] [--nilavu]...`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Megaminstall) Run(context *cmd.Context) error {
	fmt.Println("[main] starting megdc ...")

	packages := make(map[string]string)
	options := make(map[string]string)

	s := reflect.ValueOf(c).Elem()
	typ := s.Type()
	if s.Kind() == reflect.Struct {
		for i := 0; i < s.NumField(); i++ {
			key := s.Field(i)
			value := s.FieldByName(typ.Field(i).Name)
			switch key.Interface().(type) {
			case bool:
				if value.Bool() {
					packages[typ.Field(i).Name] = typ.Field(i).Name
				}
			case string:	
				if value.String() != "" {
					options[typ.Field(i).Name] = value.String()
				}
			}
		}
	}
	
	if handler, err := handler.NewHandler(); err != nil {
		return err
	} else {
		handler.SetTemplates(packages, options)
        err := handler.Run()
        if err != nil {
        	return err
        }
	}  
     
	// goodbye.
	return nil
}

func (c *Megaminstall) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
		c.Fs.BoolVar(&c.All, "all", false, "Install all megam packages")
		c.Fs.BoolVar(&c.All, "a", false, "Install all megam packages")
		c.Fs.BoolVar(&c.Megamnilavu, "megamnilavu", false, "Install nilavu package")
		c.Fs.BoolVar(&c.Megamnilavu, "n", false, "Install nilavu package")
		c.Fs.BoolVar(&c.Megamgateway, "megamgateway", false, "Install megam gateway package")
		c.Fs.BoolVar(&c.Megamgateway, "g", false, "Install megam gateway package")
		c.Fs.BoolVar(&c.Megamd, "megamd", false, "Install megamd package")
		c.Fs.BoolVar(&c.Megamd, "d", false, "Install megamd package")
		c.Fs.BoolVar(&c.Megamcommon, "megamcommon", false, "Install megamcommon package")
		c.Fs.BoolVar(&c.Megamcommon, "c", false, "Install megamcommon package")
		c.Fs.BoolVar(&c.Megammonitor, "megammonitor", false, "Install megam monitor package")
		c.Fs.BoolVar(&c.Megammonitor, "m", false, "Install megam monitor package")
		c.Fs.StringVar(&c.Host, "host", "", "host address for machine")
		c.Fs.StringVar(&c.Host, "h", "", "host address for machine")
		c.Fs.StringVar(&c.Username, "username", "", "username for hosted machine")
		c.Fs.StringVar(&c.Username, "u", "", "username for hosted machine")
		c.Fs.StringVar(&c.Password, "password", "", "password for hosted machine")
		c.Fs.StringVar(&c.Password, "p", "", "password for hosted machine")
		c.Fs.BoolVar(&c.Quiet, "quiet", false, "")
		c.Fs.BoolVar(&c.Quiet, "q", false, "")
	}
	return c.Fs
}



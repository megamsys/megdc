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
package handler

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/megamsys/libgo/cmd"
)

const (
	HOST     = "host"
	USERNAME = "username"
	PASSWORD = "password"
	PLATFORM = "platform"
)

type WrappedParms struct {
	Packages map[string]string
	Options  map[string]string
}

func (w *WrappedParms) String() string {
	wt := new(tabwriter.Writer)
	var b bytes.Buffer
	wt.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("Packages", "cyan", "", "") + "\n"))
	for _, v := range w.Packages {
		b.Write([]byte(v + "\n"))
	}
	b.Write([]byte("---\n"))
	for k, v := range w.Options {
		b.Write([]byte(k + "\t" + v + "\n"))
	}
	fmt.Fprintln(wt)
	wt.Flush()
	return strings.TrimSpace(b.String())
}

func NewWrap(c interface{}) *WrappedParms {
	w := WrappedParms{}
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
	w.Packages = packages
	w.Options = options
	return &w
}

func (w *WrappedParms) len() int {
	return len(w.Packages)
}

func (w *WrappedParms) Empty() bool {
	return w.len() == 0
}

func (w *WrappedParms) IfNoneAddPackages(p []string) {
	if w.Empty() {
		for i := range p {
			w.addPackage(p[i])
		}
	}
}

func (w *WrappedParms) addPackage(k string) {
	w.Packages[k] = k
}

func (w *WrappedParms) GetHost() (string, bool) {
	k, v := w.Options[HOST]
	return k, v
}

func (w *WrappedParms) GetUserName() (string, bool) {
	k, v := w.Options[USERNAME]
	return k, v
}

func (w *WrappedParms) GetPassword() (string, bool) {
	k, v := w.Options[PASSWORD]
	return k, v
}

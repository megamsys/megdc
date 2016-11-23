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
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	pp "github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/os"
	"github.com/megamsys/libmegdc/templates"
	_ "github.com/megamsys/libmegdc/templates/ubuntu"
	_ "github.com/megamsys/libmegdc/templates/debian"
	_ "github.com/megamsys/libmegdc/templates/centos"
	"github.com/tj/go-spin"
)

const Logo = `
	███╗   ███╗███████╗ ██████╗ ██████╗  ██████╗
	████╗ ████║██╔════╝██╔════╝ ██╔══██╗██╔════╝
	██╔████╔██║█████╗  ██║  ███╗██║  ██║██║
	██║╚██╔╝██║██╔══╝  ██║   ██║██║  ██║██║
	██║ ╚═╝ ██║███████╗╚██████╔╝██████╔╝╚██████╗
	╚═╝     ╚═╝╚══════╝ ╚═════╝ ╚═════╝  ╚═════╝
`

type Handler struct {
	writer    io.Writer
	templates []*templates.Template
	platform  string
}

func NewHandler(w *WrappedParms) (*Handler, error) {
	h := &Handler{}
	if os, err := supportedOS(); err != nil {
		return nil, err
	} else {
		h.platform = os
	}
	fmt.Println(w)
	h.SetTemplates(w)
	return h, nil
}

func (h *Handler) SetTemplates(w *WrappedParms) {
	for k, _ := range w.Packages {
		template := templates.NewTemplate()
		var v, ok = w.GetHost()
		if ok {
			template.Host = v
		}
		v, ok = w.GetUserName()
		if ok {
			template.UserName = v
		}
		v, ok = w.GetPassword()
		if ok {
			template.Password = v
		}
		template.Name = strings.Title(h.platform) + k
		template.Options = w.Options
		template.Maps = w.Maps
		h.templates = append(h.templates, template)
	}
}

func (h *Handler) Run(w io.Writer,inputs map[string]string) error {
	return templates.RunInTemplates(h.templates, func(t *templates.Template, _ chan *templates.Template) error {
		err := t.Run(w,inputs)
		if err != nil {
			return err
		}
		return nil
	}, nil, false)
}

func supportedOS() (string, error) {
	osh := os.HostOS()
	switch runtime.GOOS {
	case "linux":
		if osh != os.Ubuntu && osh != os.Debian && osh != os.CentOS {
			return "", fmt.Errorf("unsupported operating system: %v, we support ubuntu.", osh)
		}
	default:
		return "", fmt.Errorf("unsupported operating system: %v", runtime.GOOS)
	}
	return strings.ToLower(osh.String()), nil
}

//Show a spinner until our services start.
func FunSpin(vers string, logo string, task string) {
	fmt.Printf("%s %s", vers, logo)

	s := spin.New()
	for i := 0; i < 10; i++ {
		fmt.Printf("\r%s", fmt.Sprintf("%s %s %s", pp.Colorfy("starting", "green", "", "bold"), task, s.Next()))
		time.Sleep(3 * time.Millisecond)
	}
	fmt.Printf("\n")
}

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
	//"errors"
	"github.com/megamsys/megdc/templates"
	"io"
	"strings"
	_ "github.com/megamsys/megdc/templates/ubuntu"
	//"fmt"
)

const (
	HOST     = "host"
	USERNAME = "username"
	PASSWORD = "password"
	BRIDGENAME = "bridgename"
	NETWORK_IF = "networkif"
)

type Handler struct {
	writer    io.Writer
	templates []*templates.Template
	platform  string
}

func NewHandler() (*Handler, error) {
	h := &Handler{}

	if platform_name, err := findPlatform(); err != nil {
		return h, err
	} else {
		h.platform = platform_name
	}

	return h, nil

}

func (h *Handler) SetTemplates(packages map[string]string, options map[string]string) {
	for k, _ := range packages {
		template := templates.NewTemplate()
		for ko, vo := range options {
			if ko == HOST {
				template.Host = vo
			}
			if ko == USERNAME {
				template.UserName = vo
			}
			if ko == PASSWORD {
				template.Password = vo
			}
			if ko == BRIDGENAME {
				template.Bridgename = vo
			}
			if ko == NETWORK_IF {
				template.Networkif = vo
			}
		}
		template.Name = strings.Title(h.platform) + k
		h.templates = append(h.templates, template)
	}
}

func (h *Handler) Run() error {
	return templates.RunInTemplates(h.templates, func(t *templates.Template, _ chan *templates.Template) error {
		err := t.Run()
    if err != nil {
			return err
		}
		return nil
	}, nil, false)
}

func findPlatform() (string, error) {

	return "ubuntu", nil
}

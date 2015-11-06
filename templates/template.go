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

package templates

import (
	"github.com/dynport/urknall"
	"os"
	"sync"
	//"reflect"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

//var templates Templates

var templates Templates

var managers map[string]Templates

type Templates interface {
}

type InitializableTemplates interface {
	Run(target urknall.Target) error
}

type Template struct {
	Host     string
	UserName string
	Password string
	Name     string
}

const LOCALHOST = "localhost"

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) Run() error {
	defer urknall.OpenLogger(os.Stdout).Close()
	var target urknall.Target
	var e error
	//	uri := "ubuntu@192.168.56.10"
	//	password := "ubuntu"
	if t.Password != "" {
		target, e = urknall.NewSshTargetWithPassword(t.Host, t.Password)
	} else {
		if t.Host == LOCALHOST {
			target, e = urknall.NewLocalTarget()
		} else {
			///target, e = urknall.NewSshTarget(t.Host)
			target, e = urknall.NewLocalTarget()
		}
	}
	if e != nil {
		return e
	}
	a, err := get(t.Name)

	if err != nil {
		log.Errorf("fatal error, couldn't located the Package %s", t.Name)
		return err
	}

	templates = a

	if initializableTemplates, ok := templates.(InitializableTemplates); ok {
		return initializableTemplates.Run(target)
	}

	return errors.New(fmt.Sprintf("fatal error, couldn't located the Package %q", t.Name))
}

type callbackFunc func(*Template, chan *Template) error

type rollbackFunc func(*Template)

func RunInTemplates(templates []*Template, callback callbackFunc, rollback rollbackFunc, parallel bool) error {
	if len(templates) == 0 {
		return nil
	}
	workers := 0
	if workers == 0 {
		workers = len(templates)
	}
	step := len(templates)/workers + 1
	toRollback := make(chan *Template, len(templates))
	errors := make(chan error, len(templates))
	var wg sync.WaitGroup
	runFunc := func(start, end int) error {
		defer wg.Done()
		for i := start; i < end; i++ {
			err := callback(templates[i], toRollback)
			if err != nil {
				errors <- err
				return err
			}
		}
		return nil
	}
	for i := 0; i < len(templates); i++ {
		//for i := 0; i < len(templates); i += step {
		end := i + step
		if end > len(templates) {
			end = len(templates)
		}
		wg.Add(1)
		if parallel {
			go runFunc(i, end)
		} else {
			err := runFunc(i, end)
			if err != nil {
				//break
			}
		}
	}
	wg.Wait()
	close(errors)
	close(toRollback)
	if err := <-errors; err != nil {
		if rollback != nil {
			for c := range toRollback {
				rollback(c)
			}
		}
		return err
	}
	return nil
}

// Get gets the named provisioner from the registry.
func get(name string) (Templates, error) {
	p, ok := managers[name]
	if !ok {
		return nil, fmt.Errorf("unknown template: %q", name)
	}
	return p, nil
}

// Manager returns the current configured manager, as defined in the
// configuration file.
func manager(managerName string) Templates {
	if _, ok := managers[managerName]; !ok {
		managerName = "nop"
	}
	return managers[managerName]
}

// Register registers a new repository manager, that can be later configured
// and used.
func Register(name string, manager Templates) {
	if managers == nil {
		managers = make(map[string]Templates)
	}
	managers[name] = manager
}

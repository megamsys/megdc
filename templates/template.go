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
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/urknall"
)

const LOCALHOST = "localhost"

var runnables map[string]TemplateRunnable

type TemplateRunnable interface {
	Options(t *Template)
	Run(target urknall.Target) error
}

type Template struct {
	Name     string
	Host     string
	UserName string
	Password string
	Options  map[string]string
	Maps map[string][]string
}

func NewTemplate() *Template {
	return &Template{}
}

func (t *Template) Run() error {
	defer urknall.OpenLogger(os.Stdout).Close()
	var target urknall.Target
	var err error
	if t.Password != "" {
		target, err = urknall.NewSshTargetWithPassword(t.UserName+"@"+t.Host, t.Password)


	} else {
		if len(strings.TrimSpace(t.Host)) <= 0 || t.Host == LOCALHOST {
			target, err = urknall.NewLocalTarget()
		} else {
			target, err = urknall.NewSshTarget(t.UserName+"@"+t.Host) //this is with sshkey
		}
	}
	if err != nil {
		return err
	}

	runner, err := get(t.Name)

	if err != nil {
		log.Errorf("fatal error, couldn't locate the package %s", t.Name)
		return err
	}

	if initializeRunner, ok := runner.(TemplateRunnable); ok {
		initializeRunner.Options(t)
		return initializeRunner.Run(target)
	}
	return errors.New(fmt.Sprintf("fatal error, couldn't locate the package %q", t.Name))
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
func get(name string) (TemplateRunnable, error) {
	p, ok := runnables[name]
	if !ok {
		return nil, fmt.Errorf("unknown template: %q", name)
	}
	return p, nil
}

// Register registers a new repository manager, that can be later configured
// and used.
func Register(name string, runnable TemplateRunnable) {
	if runnables == nil {
		runnables = make(map[string]TemplateRunnable)
	}
	runnables[name] = runnable
}

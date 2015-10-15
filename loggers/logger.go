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
package loggers

import (
	//	"errors"
		"fmt"
		"time"
	)

var managers map[string]InitializableLogger

// LoggerManager represents a manager of application Loggers.
type InitializableLogger interface {
	Notify(name string, logs []Boxlog) error
}

type Logger interface {
}

// Boxlog represents a log entry.
type Boxlog struct {
	Date    time.Time
	Message string
	Source  string
	Name    string
}

// Get gets the named Logger from the registry.
func Get(name string) (Logger, error) {
	p, ok := managers[name]
	if !ok {
		return nil, fmt.Errorf("unknown Logger: %q", name)
	}
	return p, nil
}

// Manager returns the current configured manager, as defined in the
// configuration file.
func Manager(managerName string) InitializableLogger {
	if _, ok := managers[managerName]; !ok {
		managerName = "nop"
	}
	return managers[managerName]
}

// Register registers a new Logger manager, that can be later configured
// and used.
func Register(name string, manager InitializableLogger) {
	if managers == nil {
		managers = make(map[string]InitializableLogger)
	}
	managers[name] = manager
}
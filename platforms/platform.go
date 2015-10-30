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
package platforms

import (
//	"fmt"
	"errors"
)

var managers map[string]InitializablePlatform

type Platform interface {
}

func Read() (string, error) {
	running_platform_name, err := reader()

	if err != nil {
		return "", err
	}

	return running_platform_name, nil
}

// PlatformManager represents a manager of supported Platforms.
type InitializablePlatform interface {
	GetInstallTemplates(option string) error
	GetUninstallTemplates(option string) error
}

func Get(name string) (Platform, error) {
	p, ok := managers[name]
	if !ok {
		return nil, errors.New("Unsupported platform")
	}
	return p, nil
}

// Manager returns the current configured manager, as defined in the
// configuration file.
func Manager(managerName string) InitializablePlatform {
	if _, ok := managers[managerName]; !ok {
		managerName = "nop"
	}
	return managers[managerName]
}

// Register registers a new Platform manager, that can be later configured
// and used.
func Register(name string, manager InitializablePlatform) {
	if managers == nil {
		managers = make(map[string]InitializablePlatform)
	}
	managers[name] = manager
}

func reader() (string, error) {

	return "ubuntu", nil
}

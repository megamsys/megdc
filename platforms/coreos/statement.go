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
package coreos

import (
	"github.com/megamsys/megdc/platforms"
)

func init() {
	platforms.Register("coreos", coreosPlatform{})
}

type coreosPlatform struct {
	
}

func (u coreosPlatform) GetInstallTemplates(option string) error {
	
	return nil
}

func (u coreosPlatform) GetUninstallTemplates(option string) error {
	
	return nil
}
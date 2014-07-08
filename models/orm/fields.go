/*
** Copyright [2012-2013] [Megam Systems]
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

package orm

import (
    "github.com/megamsys/cloudinabox/modules/utils"
    "time"
)

const layout = "Jan 2, 2006 at 3:04pm (MST)"

type Users struct {
	Id int64
	Username string
	Apikey string
	Authenticated bool
	Created string
}

func NewUser(user *utils.User) Users {
	time := time.Now()
    return Users{
        Username:        user.Username,
        Apikey:          user.Api_key,
        Authenticated:   true,
        Created:         time.Format(layout),
    }
}

type Servers struct {
	Id              int64
	Name            string
	Install         bool
	InstallDate     string
	UpdateDate      string
}

func NewServer(serverName string) Servers {
	time := time.Now()
	return Servers{
		Name:   serverName,
		Install: true,
		InstallDate: time.Format(layout),
		UpdateDate: "",
	}
}



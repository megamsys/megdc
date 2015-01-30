/*
** Copyright [2012-2014] [Megam Systems]
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
	Id              int64      `db:"Id"`
	Name            string     `db:"Name"`
	Install         bool       `db:"Install"`
	IP              string     `db:"IP"`
	Stype           string     `db:"Stype"`
	HostName        string     `db:"HostName"`
	InstallDate     string     `db:"InstallDate"`
	UpdateDate      string     `db:"UpdateDate"`
}

type Nodes struct {
	Id              int64      `db:"Id"`
	Install         bool       `db:"Install"`
	IP              string     `db:"IP"`
	HostName        string     `db:"HostName"`
	InstallDate     string     `db:"InstallDate"`
	UpdateDate      string     `db:"UpdateDate"`
}

func NewServer(serverName string, ip string, stype string, hostname string) Servers {
	time := time.Now()
	return Servers{
		Name:   serverName,
		Install: false,
		IP: ip,
		Stype: stype,
		HostName: hostname,
		InstallDate: time.Format(layout),
		UpdateDate: time.Format(layout),
	}
}


func NewNode(ip string, hostname string) Nodes {
	time := time.Now()
	return Nodes{
		Install: true,
		IP: ip,
		HostName: hostname,
		InstallDate: time.Format(layout),
		UpdateDate: time.Format(layout),
	}
}

type HAServers struct {
	Id            int64    `db:"id"`
	NodeIP1       string   `db:"nodeip1"`
	NodeHost1     string   `db:"nodehost1"`
	NodeDisk1     string   `db:"nodedisk1"`
	NodeIP2       string   `db:"nodeip2"`
	NodeHost2     string   `db:"nodehost2"`  
	NodeDisk2     string   `db:"nodedisk2"`
}


type Storages struct {
	Id           int64    `db:"Id"`
	IP       	 string   `db:"ip"`
	Storage1     string   `db:"storage1"`
	Storage2     string   `db:"storage2"`
	Storage3     string   `db:"storage3"`
}



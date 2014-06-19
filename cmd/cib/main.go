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
package main

import (
	"fmt"
	"github.com/megamsys/cloudinabox/cmd"
	"github.com/tsuru/config"
	"log"
	"os"
	"path/filepath"
)

const (
	version = "0.3.0"
	header  = "Supported-Gulp"
)

const defaultConfigPath = "conf/cib.conf"

func buildManager(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version, header)
	m.Register(&GulpStart{m, nil, false}) //start the gulpd daemon
	m.Register(&GulpStop{})               //stop  the gulpd daemon
	m.Register(&GulpUpdate{})             //stop  the gulpd daemon
	m.Register(&AppStart{})               //sudo service <appname> start
	m.Register(&AppStop{})                //sudo service <appname> stop
	/*m.Register(&AppRestart{}) //sudo service <apppname> restart
	m.Register(&AppBuild{})   //git fetch -q
	m.Register(&gulp.AppMaintain{})//sudo service nginx maintain ?
	m.Register(&gulp.SSLAdd{})     //download node_name.pub, crt from S3, mk ssl_template, cp to sites_available, ln to sites_enabled. && AppRestart
	m.Register(&gulp.SSLRemove{})  //rm node_name.pub, crt, mk regular non_ssl_template, cp to sites_available, ln to sites_enabled. && AppRestart
	m.Register(&gulp.MeterStop{})  //sudo service gmond start
	m.Register(&gulp.MeterStart{}) //sudo service gmond stop
	m.Register(&gulp.LogStart{})   //sudo service beaver start
	m.Register(&gulp.LogStop{})    //sudo service beaver stop
	m.Register(&gulp.EnvGet{})     //ENV['JMP_UP_PATH'] '~/sofware/kangaroo'
	m.Register(&gulp.EnvSet{})     //ENV['JMP_UP_PATH'] = '~/software/kangaroo'
	m.Register(&gulp.EnvUnset{})   //ENV['JMP_UP_PATH'] = blank
	m.Register(&KeyAdd{})          //add the id_rsa/pub
	m.Register(&KeyRemove{})       //remove the id_rsa/pub
	m.Register(gulp.ServiceList{}) //ps -ef
	*/
	return m
}

func main() {
	p, _ := filepath.Abs(defaultConfigPath)
	log.Println(fmt.Errorf("Conf: %s", p))
	config.ReadConfigFile(defaultConfigPath)
	name := cmd.ExtractProgramName(os.Args[0])
	manager := buildManager(name)
	manager.Run(os.Args[1:])
}

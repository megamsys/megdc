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

package ubuntu

import (
	"os"
	"github.com/megamsys/megdc/templates"
	"github.com/megamsys/urknall"

	//"github.com/megamsys/libgo/cmd"
)

const (
		MasterIp    = "MasterIp"
	SparkHome = "/var/lib/meglytics/"
  Sparketc = "spark"
  SparkVersion = "SparkVersion"
	)

var ubuntumesosslaveinstall *UbuntuMesosSlaveInstall

func init() {
	ubuntumesosslaveinstall = &UbuntuMesosSlaveInstall{}
	templates.Register("UbuntuMesosSlaveInstall", ubuntumesosslaveinstall)
}

type UbuntuMesosSlaveInstall struct {

	masterip    string
  sparkversion string
}

func (tpl *UbuntuMesosSlaveInstall) Options(t *templates.Template) {

	if masterip, ok := t.Options[MasterIp]; ok {
		tpl.masterip = masterip
}
if sparkversion, ok := t.Options[SparkVersion]; ok {
  tpl.sparkversion = sparkversion
}

}
func (tpl *UbuntuMesosSlaveInstall) Render(p urknall.Package) {
	p.AddTemplate("messosslave", &UbuntuMesosSlaveInstallTemplate{

masterip		:    tpl.masterip,
sparkversion :  tpl.sparkversion,

	})
}

func (tpl *UbuntuMesosSlaveInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMesosSlaveInstall{

		masterip:    tpl.masterip,
    sparkversion: tpl.sparkversion,

	})
}

type UbuntuMesosSlaveInstallTemplate struct {

	masterip    string
  sparkversion string

}

func (m *UbuntuMesosSlaveInstallTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()
	ip := IP("")
	MasterIp := m.masterip
  SparkVersion := m.sparkversion
  sparkhome := SparkHome
  spark := Sparketc

	pkg.AddCommands("meglyticsdir",
  Shell("cd /var/lib;mkdir meglytics"),
  )
	
  pkg.AddCommands("MesospherRepo",
	 Shell("apt-key adv --keyserver keyserver.ubuntu.com --recv E56151BF"),
   Shell("DISTRO=$(lsb_release -is | tr '[:upper:]' '[:lower:]')"),
   Shell("CODENAME=$(lsb_release -cs)"),
   Shell("echo 'deb http://repos.mesosphere.io/ubuntu trusty main'|sudo tee /etc/apt/sources.list.d/mesosphere.list"),
	)
	pkg.AddCommands("update",
		Shell("apt-get -y update"),
	)
	pkg.AddCommands("installmesos",
		Shell("apt-get -y install mesos"),
	)


  pkg.AddCommands("zookeeper",
         Shell("echo manual | sudo tee /etc/init/zookeeper.override"),
         Shell("apt-get -y remove --purge zookeeper"),
      )
  pkg.AddCommands("slaveip",
    	Shell("echo "+ip+" | sudo tee /etc/mesos-slave/ip"),
      Shell("echo "+ip+" | sudo tee /etc/mesos-slave/hostname"),
    )
    pkg.AddCommands("masterip",
      	Shell("echo zk://"+MasterIp+":2181/mesos | sudo tee /etc/mesos/zk"),
      )

	pkg.AddCommands("etchost",
		Shell("echo '"+ip+" "+host+"' >> /etc/hosts"),
    Shell("echo '"+ip+" "+spark+"' >> /etc/hosts"),
	)
  pkg.AddCommands("slavestart",
      Shell("service mesos-slave restart"),
    )
    pkg.AddCommands("sparkinstall",
        Shell("cd /home;wget http://www.eu.apache.org/dist/spark/spark-"+SparkVersion+"/spark-"+SparkVersion+".tgz "),
				Shell("cp /home/spark-"+SparkVersion+".tgz "+sparkhome+""),
        Shell(" cd "+sparkhome+";tar xvf spark-"+SparkVersion+".tgz"),
        Shell("mv "+sparkhome+"spark-"+SparkVersion+" "+sparkhome+"spark "),
      )
      pkg.AddCommands("sparkenvsh",
          Shell("echo export SPARK_EXECUTOR_URI=/home/spark-"+SparkVersion+".tgz >> /var/lib/meglytics/spark/conf/spark-env.sh.template"),
         Shell(" echo export MESOS_NATIVE_LIBRARY=/usr/local/lib/libmesos.so >> /var/lib/meglytics/spark/conf/spark-env.sh.template"),
        )
        pkg.AddCommands("sparkdefaults",
         Shell("echo spark.master mesos://"+MasterIp+":5050 >> /var/lib/meglytics/spark/conf/spark-defaults.conf.template"),
         Shell("echo spark.executor.uri /home/spark-"+SparkVersion+".tgz >> /var/lib/meglytics/spark/conf/spark-defaults.conf.template"),
        )
}

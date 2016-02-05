package ubuntu

import (

	"os"
   "fmt"
	"github.com/megamsys/megdc/templates"
	"github.com/megamsys/urknall"
	//"github.com/megamsys/libgo/cmd"
)

const (
		Sparkhome= "/var/lib/meglytics/"
		spark = "spark"
    Sparkvrs = "Sparkvrs"
 MEGAMCONF = `
 # Template for a Spark Job Server configuration file
# When deployed these settings are loaded when job server starts
#
# Spark Cluster / Job Server configuration
spark {
  # spark.master will be passed to each job's JobContext
    master="%s"
  # master = "mesos://vm28-hulk-pub:5050"
  # master = "yarn-client"

  # Default # of CPUs for jobs to use for Spark standalone cluster
  job-number-cpus = 4

  jobserver {
    port = 8090
    jar-store-rootdir = /tmp/jobserver/jars

    context-per-jvm = "%s"

    jobdao = spark.jobserver.io.JobFileDAO

    filedao {
      rootdir = /tmp/spark-job-server/filedao/data
    }
  }

  # predefined Spark contexts
  # contexts {
  #   my-low-latency-context {
  #     num-cpu-cores = 1           # Number of cores to allocate.  Required.
  #     memory-per-node = 512m         # Executor memory per node, -Xmx style eg 512m, 1G, etc.
  #   }
  #   # define additional contexts here
  # }

  # universal context configuration.  These settings can be overridden, see README.md
    context-settings {

    num-cpu-cores = 2           # Number of cores to allocate.  Required.
    memory-per-node = 512m         # Executor memory per node, -Xmx style eg 512m, #1G, etc.

    # in case spark distribution should be accessed from HDFS (as opposed to being installed on every mesos slave)
    # spark.executor.uri = "hdfs://namenode:8020/apps/spark/spark.tgz"
     spark.mesos.executor.uri = "%s"
    # uris of jars to be loaded into the classpath for this context. Uris is a string list, or a string separated by commas ','
    # dependent-jar-uris = ["file:///some/path/present/in/each/mesos/slave/somepackage.jar"]

    # If you wish to pass any settings directly to the sparkConf as-is, add them here in passthrough,
    # such as hadoop connection settings that don't use the "spark." prefix
    passthrough {
      #es.nodes = "192.1.1.1"
    }
  }

  # This needs to match SPARK_HOME for cluster SparkContexts to be created successfully
   home = "%s"
}

# Note that you can use this file to define settings not only for job server,
# but for your Spark jobs as well.  Spark job configuration merges with this configuration file as defaults.
`
MEGAMSH = `

# Environment and deploy file
# For use with bin/server_deploy, bin/server_package etc.
DEPLOY_HOSTS="%s"
APP_USER="%s"
APP_GROUP="%s"
# optional SSH Key to login to deploy server
#SSH_KEY=/path/to/keyfile.pem
INSTALL_DIR=%s
LOG_DIR=/var/log/job-server
PIDFILE=spark-jobserver.pid
JOBSERVER_MEMORY=1G
SPARK_VERSION=%s
SPARK_HOME=%s
SPARK_CONF_DIR=$SPARK_HOME/conf
# Only needed for Mesos deploys
SPARK_EXECUTOR_URI=%s
# Only needed for YARN running outside of the cluster
# You will need to COPY these files from your cluster to the remote machine
# Normally these are kept on the cluster in /etc/hadoop/conf
# YARN_CONF_DIR=/pathToRemoteConf/conf
# HADOOP_CONF_DIR=/pathToRemoteConf/conf
#
# Also optional: extra JVM args for spark-submit
# export SPARK_SUBMIT_OPTS+="-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=5433"
SCALA_VERSION=2.10.4 # or 2.11.6
`
	)

var ubuntumesosmasterinstall *UbuntuMesosMasterInstall

func init() {
	ubuntumesosmasterinstall = &UbuntuMesosMasterInstall{}
	templates.Register("UbuntuMesosMasterInstall", ubuntumesosmasterinstall)
}

type UbuntuMesosMasterInstall struct {
	sparkvrs string

}

func (tpl *UbuntuMesosMasterInstall) Options(t *templates.Template) {
	if sparkvrs, ok := t.Options[Sparkvrs]; ok {
	  tpl.sparkvrs = sparkvrs
	}


}

func (tpl *UbuntuMesosMasterInstall) Render(p urknall.Package) {
	p.AddTemplate("mesosmaster", &UbuntuMesosMasterInstallTemplate{
		sparkvrs :  tpl.sparkvrs,

	})
}

func (tpl *UbuntuMesosMasterInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMesosMasterInstall{
		sparkvrs: tpl.sparkvrs,

	})
}

type UbuntuMesosMasterInstallTemplate struct {
	sparkvrs string

}

func (m *UbuntuMesosMasterInstallTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()
	ip := IP("")
  sparkhome := Sparkhome
  mesos := "mesos://"+ip+":5050"
	Sparkvrs := m.sparkvrs
		job := ""+sparkhome+"spark-"+Sparkvrs+""
		executeuri := "/home/spark-"+Sparkvrs+".tgz"
		installdir := ""+sparkhome+"jobserver"

   pkg.AddCommands("meglyticsdir",
   Shell("cd /var/lib;mkdir meglytics1"),
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


   pkg.AddCommands("ipofmaster",
 		 Shell("echo "+ip+" | sudo tee /etc/mesos-master/ip"),
 	 )

   pkg.AddCommands("zookeeper",
     Shell("echo zk://"+ip+":2181/mesos | sudo tee /etc/mesos/zk"),
   )

   pkg.AddCommands("hostname",
     Shell("echo "+ip+" | sudo tee /etc/mesos-master/hostname"),
   )

	pkg.AddCommands("etchost",
		Shell("echo '"+ip+" "+host+"' >> /etc/hosts"),
		Shell("echo '"+ip+" "+spark+"' >> /etc/hosts"),
	)

	pkg.AddCommands("zookeeperstart",
		Shell("service zookeeper restart"),
	)

  pkg.AddCommands("masterstart",
    Shell("service mesos-master restart"),
  )
	pkg.AddCommands("sparkinstall",
			Shell("cd /home;wget http://www.eu.apache.org/dist/spark/spark-"+Sparkvrs+"/spark-"+Sparkvrs+".tgz "),
			Shell("cp /home/spark-"+Sparkvrs+".tgz "+sparkhome+""),
			Shell(" cd "+sparkhome+";tar xvf spark-"+Sparkvrs+".tgz"),
			Shell("mv "+sparkhome+"spark-"+Sparkvrs+" "+sparkhome+"spark "),
		)
		pkg.AddCommands("sparkenvsh",
				Shell("echo export SPARK_EXECUTOR_URI=/home/spark-"+Sparkvrs+".tgz >> /var/lib/meglytics/spark/conf/spark-env.sh.template"),
			 Shell(" echo export MESOS_NATIVE_LIBRARY=/usr/local/lib/libmesos.so >> /var/lib/meglytics/spark/conf/spark-env.sh.template"),
			)
			pkg.AddCommands("sparkdefaults",
			 Shell("echo spark.master mesos://"+ip+":5050 >> /var/lib/meglytics/spark/conf/spark-defaults.conf.template"),
			 Shell("echo spark.executor.uri /home/spark-"+Sparkvrs+".tgz >> /var/lib/meglytics/spark/conf/spark-defaults.conf.template"),
			)

  pkg.AddCommands("sparkjobserverclone",
    Shell("cd "+sparkhome+";git clone https://github.com/spark-jobserver/spark-jobserver.git"),
  )
  pkg.AddCommands("localconfchange",
    Shell("cp "+sparkhome+"spark-jobserver/config/local.conf.template "+sparkhome+"spark-jobserver/config/megam.conf"),
  )
  pkg.AddCommands("localshchange",
    Shell("cp "+sparkhome+"spark-jobserver/config/local.sh.template   "+sparkhome+"spark-jobserver/config/megam.sh"),
  )

   pkg.AddCommands("conf",

WriteFile(sparkhome+"spark-jobserver/config/megam.conf", fmt.Sprintf(MEGAMCONF, mesos, "false", job, job), "root", 0755),
)

pkg.AddCommands("export",
 Shell("export SPARK_HOME='/var/lib/meglytics/spark-"+Sparkvrs+"'"),
)

  pkg.AddCommands("sh",
    WriteFile(sparkhome+"spark-jobserver/config/megam.sh", fmt.Sprintf(MEGAMSH, ip, "root", "root", installdir, Sparkvrs, job, executeuri), "root", 0755),
)
pkg.AddCommands("run",
Shell("cd "+sparkhome+"spark-jobserver/bin;./server_deploy.sh megam"),
)
pkg.AddCommands("cadvisor",
Shell("stop cadvisor"),
)
pkg.AddCommands("job",
	Shell("cd "+sparkhome+"jobserver;./server_start.sh"),
)
}

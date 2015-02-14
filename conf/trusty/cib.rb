require 'json'

#Packages list
@packages = {
"megam" => {"megamcommon" => "false", "megamcib" => "false", "megamcibn" => "false", "megamnilavu" => "false", "megamsnowflake" => "false", "megamgateway" => "false", "megamd" => "false", "chef-server" => "false", "megamanalytics" => "false", "megamvarai" => "false", "megammonitor" => "false", "riak" => "false", "rabbitmq-server" => "false", "nodejs" => "false", "sqlite3" => "false", "ruby2.0" => "false", "openjdk-7-jdk" => "false"},
"cobbler" => {"cobbler" => "false", "dnsmasq" => "false", "apache2" => "false", "debmirror" => "false"},
"opennebula" => {"opennebula" => "false", "opennebula-sunstone" => "false"},
"opennebula_host" => {"opennebula-node" => "false", "qemu-kvm" => "false"},
"ceph" => {"ceph-deploy" => "false", "ceph-common ceph-mds" => "false"},
"drbd" => {"drbd8-utils" => "false", "linux-image-extra-virtual" => "false", "pacemaker" => "false", "heartbeat" => "false"}
}

#Package installation check for debian
def pkg_check(pkg_array)
        pkg_array.each do |pkg|
          @packages["#{pkg}"].each_key do |k|
                dpkg_res = `dpkg -s #{k}`
                if "#{dpkg_res}".include? "Status: install ok installed"
                        @packages["#{pkg}"]["#{k}"] = "true"
                else
                        puts "Package #{k} Not Installed"
                end
          end
        end
@packages.select{|k, _| pkg_array.include?(k)}
end

#pkg_check(["megam", "ceph"])



#Some of the packages don't have services
#megamcommon, megamsnowflake(but snowflake),megammonitor(but ganglia-monitor), nodejs, sqlite3, ruby2.0, openjdk-7-jdk, debmirror
@services = {
"megam" => {"megamcommon" => "false", "megamcib" => "false", "megamcibn" => "false", "megamnilavu" => "false", "snowflake" => "false", "megamgateway" => "false", "megamd" => "false", "chef-server-ctl" => "false", "megamanalytics" => "false", "megamvarai" => "false", "ganglia-monitor" => "false", "riak" => "false", "rabbitmq-server" => "false", "nodejs" => "false", "sqlite3" => "false", "ruby2.0" => "false", "openjdk-7-jdk" => "false"},
"cobbler" => {"cobbler" => "false", "dnsmasq" => "false", "apache2" => "false"},
"opennebula" => {"opennebula" => "false", "opennebula-sunstone" => "false"},
"opennebula_host" => {"opennebula-node" => "false", "qemu-kvm" => "false"},
"ceph" => {"ceph-all" => "false", "ceph_health" => "false"},
"drbd" => {"drbd" => "false", "pacemaker" => "false", "heartbeat" => "false", "crm" => "false"}
}

def service_check(service_array)
        service_array.each do |ser|
          @services["#{ser}"].each_key do |k|
                if `sudo service #{k} status`.include? "running"
                         @services["#{ser}"]["#{k}"] = "true"
                end
                if `sudo service #{k} status`.include? "not"
                         @services["#{ser}"]["#{k}"] = "false"
                end
                if ("#{k}" == "riak") && (`sudo riak ping`.include? "pong")
                         @services["#{ser}"]["#{k}"] = "true"
                end
                if ("#{k}" == "chef-server-ctl") && (`sudo chef-server-ctl status`.include? "run")
                         @services["#{ser}"]["#{k}"] = "true"
                end
                if ("#{k}" == "drbd") && (`sudo service #{k} status`.include? "drbd driver loaded OK")
                         @services["#{ser}"]["#{k}"] = "true"
                end
                if ("#{k}" == "ceph_health") && (`ceph health`.include? "HEALTH_OK")
                        @services["#{ser}"]["ceph_health"] = "true"
                end
                if ("#{k}" == "crm") && (`crm status`.include? "Started")
                        @services["#{ser}"]["crm"] = "true"
                end
          end
        end
@services.select{|k, _| service_array.include?(k)}
end

#service_check(["megam", "ceph"])


def check_cib(array)
        cib_check = {"packages" => {}, "services" => {}}
        pkg_hash = pkg_check(array)
        service_hash = service_check(array)
        cib_check["packages"] = pkg_hash
        cib_check["services"] = service_hash
cib_check.to_json
end


cib_json = check_cib(ARGV)
puts cib_json



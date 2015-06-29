require 'json'

#Packages list
@packages = {
"megam" => {"megamcommon" => "true", "megamcib" => "true", "megamcibn" => "true", "megamnilavu" => "true", "megamsnowflake" => "true", "megamgateway" => "true", "megamd" => "true", "chef-server" => "true", "megamanalytics" => "true", "megamvarai" => "true", "megammonitor" => "true", "riak" => "true", "rabbitmq-server" => "true", "nodejs" => "true", "sqlite3" => "true", "ruby2.0" => "true", "openjdk-7-jdk" => "true"},
"cobbler" => {"cobbler" => "true", "dnsmasq" => "true", "apache2" => "true", "debmirror" => "true"},
"opennebula" => {"opennebula" => "true", "opennebula-sunstone" => "true"},
"opennebulahost" => {"opennebula-node" => "true", "qemu-kvm" => "true"},
"ceph" => {"ceph-deploy" => "true", "ceph-common ceph-mds" => "true"},
"drbd" => {"drbd8-utils" => "true", "linux-image-extra-virtual" => "true", "pacemaker" => "true", "heartbeat" => "true"}
}

#Package installation check for debian
def pkg_check(pkg_array)
        pkg_array.each do |pkg|
          @packages["#{pkg}"].each_key do |k|

                dpkg_res = `dpkg -s #{k} >/dev/null 2>&1 && { printf "success"; } || { printf "fail";}`
                if "#{dpkg_res}".include? "success"
                        @packages["#{pkg}"]["#{k}"] = "true"
                end
          end
        end
@packages.select{|k, _| pkg_array.include?(k)}
end

#pkg_check(["megam", "ceph"])

#puts @packages

#Some of the packages don't have services
#megamcommon, megamsnowflake(but snowflake),megammonitor(but ganglia-monitor), nodejs, sqlite3, ruby2.0, openjdk-7-jdk, debmirror
@services = {
"megam" => {"megamcommon" => "true", "megamcib" => "true", "megamcibn" => "true", "megamnilavu" => "true", "snowflake" => "true", "megamgateway" => "true", "megamd" => "true", "chef-server-ctl" => "true", "megamanalytics" => "true", "megamvarai" => "true", "ganglia-monitor" => "true", "riak" => "true", "rabbitmq-server" => "true", "nodejs" => "true", "sqlite3" => "true", "ruby2.0" => "true", "openjdk-7-jdk" => "true"},
"cobbler" => {"cobbler" => "true", "dnsmasq" => "true", "apache2" => "true"},
"opennebula" => {"opennebula" => "true", "opennebula-sunstone" => "true"},
"opennebulahost" => {"opennebula-node" => "true", "qemu-kvm" => "true"},
"ceph" => {"ceph-all" => "true", "ceph_health" => "true"},
"drbd" => {"drbd" => "true", "pacemaker" => "true", "heartbeat" => "true", "crm" => "true"}
}

def service_check(service_array)
        service_array.each do |ser|
          @services["#{ser}"].each_key do |k|
                if `sudo service #{k} status >/dev/null 2>&1 && { printf "success"; } || { printf "fail";}`.include? "success"
                         @services["#{ser}"]["#{k}"] = "true"
                end
                if ("#{k}" == "chef-server-ctl") && (`sudo chef-server-ctl status >/dev/null 2>&1 && { printf "success"; } || { printf "fail";}`.include? "success")
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

#puts @services


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


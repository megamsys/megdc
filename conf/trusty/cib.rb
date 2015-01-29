require 'json'
require 'open3'

#Packages list
@packages = {
"megam" => {"megamcommon" => "false", "megamcib" => "false", "megamcibn" => "false", "megamnilavu" => "false", "megamsnowflake" => "false", "megamgateway" => "false", "megamd" => "false", "megamchefnative" => "false", "megamanalytics" => "false", "megamdesigner" => "false", "megammonitor" => "false", "riak" => "false", "rabbitmq-server" => "false", "nodejs" => "false", "sqlite3" => "false", "ruby2.0" => "false", "openjdk-7-jdk" => "false"},
"cobbler" => {"cobbler" => "false", "dnsmasq" => "false", "apache2" => "false", "debmirror" => "false"},
"opennebula" => {"opennebula" => "false", "opennebula-sunstone" => "false"},
"opennebula_host" => {"opennebula-node" => "false", "qemu-kvm" => "false"},
"ceph" => {"ceph-deploy" => "false", "ceph-common" => "flase", "ceph-mds" => "false"},
"drbd" => {"drbd8-utils" => "false", "linux-image-extra-virtual" => "false", "pacemaker" => "false", "heartbeat" => "false"}
}

#Package installation check for debian
def pkg_check(pkg_array)
        pkg_array.each do |pkg|
          @packages["#{pkg}"].each_key do |k|
             stdin, stdout, stderr = Open3.popen3("dpkg -s #{k}")
                if stdout.readlines.count > 0
                        @packages["#{pkg}"]["#{k}"] = "true"
                end
          end
        end
return @packages.select{|k, _| pkg_array.include?(k)}
end

#pkg_check(["megam", "ceph"])

#puts @packages

@services = {
"megam" => {"megamcommon" => "false", "megamcib" => "false", "megamcibn" => "false", "megamnilavu" => "false", "snowflake" => "false", "megamgateway" => "false", "megamd" => "false", "megamchefnative" => "false", "megamanalytics" => "false", "megamdesigner" => "false", "megammonitor" => "false", "riak" => "false", "rabbitmq-server" => "false", "nodejs" => "false", "sqlite3" => "false", "ruby2.0" => "false", "openjdk-7-jdk" => "false"},
"cobbler" => {"cobbler" => "false", "dnsmasq" => "false", "apache2" => "false", "debmirror" => "false"},
"opennebula" => {"opennebula" => "false", "opennebula-sunstone" => "false"},
"opennebula_host" => {"opennebula-node" => "false", "qemu-kvm" => "false"},
"ceph" => {"ceph-all" => "false", "ceph_health" => "false"},
"drbd" => {"drbd" => "false", "pacemaker" => "false", "heartbeat" => "false", "crm" => "false"}
}

def service_check(service_array)
        service_array.each do |ser|
          @services["#{ser}"].each_key do |k|
                stdin, stdout, stderr = Open3.popen3("sudo service #{k} status")
                ser_out = stdout.readlines.map { |i| "" + i.to_s + "" }.join("")
                if ser_out.include? "running"
                         @services["#{ser}"]["#{k}"] = "true"
                end
                if ser_out.include? "not"
                         @services["#{ser}"]["#{k}"] = "false"
                end
               if ("#{k}" == "drbd") && ("#{@packages["drbd"]["drbd8-utils"]}" == "true") && (ser_out.include? "drbd driver loaded OK")
                         @services["#{ser}"]["#{k}"] = "true"
                end
                if ("#{k}" == "ceph_health") && ("#{@packages["ceph"]["ceph-common"]}" == "true") && (`ceph health`.include? "HEALTH_OK")
                        @services["#{ser}"]["ceph_health"] = "true"
                end
                if ("#{k}" == "crm") && ("#{@packages["drbd"]["pacemaker"]}" == "true") && (`sudo crm status`.include? "Started")
                        @services["#{ser}"]["crm"] = "true"
                end
          end
        end
return @services.select{|k, _| service_array.include?(k)}
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


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

package debian

import (
	"fmt"
	"os"
	"github.com/megamsys/megdc/templates"
	u "github.com/megamsys/megdc/templates/ubuntu"
	"github.com/megamsys/urknall"
	//"github.com/megamsys/libgo/cmd"
)

const (
  Ceph_conf = `
[client.radosgw.admin]
host = %s
keyring = /etc/ceph/ceph.client.radosgw.keyring
rgw socket path =   /var/run/ceph/ceph.radosgw.admin.fastcgi.sock
log file = /var/log/radosgw/client.radosgw.admin.log
`
S3gw =`
#!/bin/sh
exec /usr/bin/radosgw -c /etc/ceph/ceph.conf -n client.radosgw.admin
`
Rgw_conf =`
FastCgiExternalServer /var/www/html/s3gw.fcgi -socket /var/run/ceph/ceph.radosgw.admin.fastcgi.sock

<VirtualHost *:80>
ServerName localhost
DocumentRoot /var/www/html

ErrorLog /var/log/apache2/rgw_error.log
CustomLog /var/log/apache2/rgw_access.log combined

# LogLevel debug

RewriteEngine On

RewriteRule ^/([a-zA-Z0-9-_.]*)([/]?.*) /s3gw.fcgi?page=$1&params=$2&%{QUERY_STRING} [E=HTTP_AUTHORIZATION:%{HTTP:Authorization},L]

<IfModule mod_fastcgi.c>
    <Directory /var/www/html>
    Options +ExecCGI
    AllowOverride All
    SetHandler fastcgi-script
    Order allow,deny
    Allow from all
    AuthBasicAuthoritative Off
    </Directory>
</IfModule>

AllowEncodedSlashes On
ServerSignature Off

</VirtualHost>
`
)

var debiancephgateway *DebianCephGateway

func init() {
	debiancephgateway = &DebianCephGateway{}
	templates.Register("DebianCephGateway", debiancephgateway)
}

type DebianCephGateway struct {
}

func (tpl *DebianCephGateway) Options(t *templates.Template) {
}

func (tpl *DebianCephGateway) Render(p urknall.Package) {
	p.AddTemplate("ceph", &DebianCephGatewayTemplate{
	})
}

func (tpl *DebianCephGateway) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianCephGateway{
	})
}

type DebianCephGatewayTemplate struct {
}

func (m *DebianCephGatewayTemplate) Render(pkg urknall.Package) {
host, _ := os.Hostname()
	pkg.AddCommands("apache2",
		u.UpdatePackagesOmitError(),
		u.InstallPackages("apache2 libapache2-mod-fastcgi"),
		u.Shell("echo 'ServerName "+host+"' >>/etc/apache2/apache2.conf"),
	  u.Shell("sudo a2enmod rewrite"),
    u.Shell("sudo a2enmod fastcgi"),
		u.Shell("service apache2 start"),
	)
	pkg.AddCommands("gateway-daemon",
		u.InstallPackages("radosgw"),
		u.Shell("ceph-authtool --create-keyring /etc/ceph/ceph.client.radosgw.keyring"),
		u.Shell("chmod +r /etc/ceph/ceph.client.radosgw.keyring"),
		u.Shell("ceph-authtool /etc/ceph/ceph.client.radosgw.keyring -n client.radosgw.admin --gen-key"),
		u.Shell("ceph-authtool -n client.radosgw.admin --cap osd 'allow rwx' --cap mon 'allow rwx' /etc/ceph/ceph.client.radosgw.keyring"),
		u.Shell("ceph -k /etc/ceph/ceph.client.admin.keyring auth add client.radosgw.admin -i /etc/ceph/ceph.client.radosgw.keyring"),
	)
/*	pkg.AddCommands("copy_keyring",
		u.Shell("sudo scp /etc/ceph/ceph.client.radosgw.keyring  "+ USERNAME+"@"+GATEWAY_IP+":/home/ceph"),
		u.Shell("ssh "+ USERNAME+"@"+GATEWAY_IP+" 'sudo mv ceph.client.radosgw.keyring /etc/ceph/ceph.client.radosgw.keyring'"),
	)*/
	pkg.AddCommands("create-osd-pools",
		u.Shell("ceph osd pool create .rgw.buckets 16 16"),
		u.Shell("ceph osd pool create .rgw 16 16"),
		u.Shell("ceph osd pool create .rgw.root 16 16"),
		u.Shell("ceph osd pool create .rgw.control 16 16"),
		u.Shell("ceph osd pool create .rgw.gc 16 16"),
		u.Shell("ceph osd pool create .rgw..buckets.index 16 16"),
		u.Shell("ceph osd pool create .log 16 16"),
		u.Shell("ceph osd pool create .intent-log 16 16"),
		u.Shell("ceph osd pool create .usage 16 16"),
		u.Shell("ceph osd pool create .users 16 16"),
		u.Shell("ceph osd pool create .users.email 16 16"),
		u.Shell("ceph osd pool create .users.swift 16 16"),
		u.Shell("ceph osd pool create .users.uid 16 16"),
		u.Shell("rados lspools"),
	)
  pkg.AddCommands("ceph-conf",
		u.Shell("cat >> /etc/ceph/ceph.conf <<'EOF' "+fmt.Sprintf(Ceph_conf,host)+"EOF"),
		u.Shell("ceph-deploy --overwrite-conf config pull "+host),
    u.Shell("ceph-deploy --overwrite-conf config push "+host),
	)
	/*	pkg.AddCommands("copy_keyring",
			u.Shell("sudo scp /etc/ceph/ceph.client.admin.keyring  "+ USERNAME+"@"+GATEWAY_IP+":/home/"+USERNAME),
			u.Shell("ssh "+ USERNAME+"@"+GATEWAY_IP+" 'sudo mv ceph.client.admin.keyring /etc/ceph/ceph.client.admin.keyring'"),
		)*/
	pkg.AddCommands("Cgi-wrapper",
		u.WriteFile("/var/www/html/s3gw.fcgi",S3gw,"root",0755),
		u.Mkdir("/var/lib/ceph/radosgw/ceph-radosgw.admin","",0755),
		u.Shell("sudo /etc/init.d/radosgw start"),
		u.WriteFile("/etc/apache2/sites-available/rgw.conf",Rgw_conf,"root",0644),
		u.Shell("sudo a2dissite 000-default"),
		u.Shell("sudo a2ensite rgw.conf"),
		u.Shell("sudo service apache2 restart"),
	)
}

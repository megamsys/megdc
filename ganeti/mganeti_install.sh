#!/bin/bash

#download ganeti package
#install ganeti deb package
wget http://downloads.ganeti.org/releases/2.11/ganeti-2.11.0~rc1.tar.gz
tar xvzf ganeti-2.11.0.tar.gz
cd ganeti-2.11.0

./configure --localstatedir=/var --sysconfdir=/etc --enable-symlinks 
sudo make
sudo make install
sudo mkdir /srv/ganeti/ /srv/ganeti/os /srv/ganeti/export
sudo cp doc/examples/ganeti.initd /etc/init.d/ganeti

#Setup Ganeti service

sudo update-rc.d ganeti defaults 20 80

sudo mkdir default
sudo ln -s /usr/local/share/ganeti/2.11/ ./default 



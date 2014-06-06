#!/bin/bash

lsmod | grep drbd  >> verify_out

#Verify hostname: 
cat /etc/hostname >> verify_out

#Verify if kvm exists
kvm-ok  | grep "KVM acceleration can be used"  >> verify_out

#Verify FQDN
hostname --fqdn  | grep "?ganeti.megam.com"  >> verify_out


#verify LVM
#If the line starts with UUID=xyz, this means it's a physical partition.
#If the line starst with /dev/sdaX, it also means it's a physical partition.
#The indicator for LVM would be something with /dev/mapper/xyz.
cat /etc/fstab  >> verify_out

brctl show >> verify_out






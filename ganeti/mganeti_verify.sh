#!/bin/bash

lsmod | grep drbd  >> verify_out

#Verify hostname: 
cat /etc/hostname >> verify_out

#Verify if kvm exists
kvm=`kvm-ok  | grep "KVM acceleration can be used"`
echo "$kvm" >> verify_out

if [ "$kvm" != "KVM acceleration can be used" ]; # Did the command work?
then # Fail
     echo "An error occured. Error: KVM doesn't exist"     
     exit 1
fi

#Verify FQDN
hostname=`hostname --fqdn`
echo "$hostname" >> verify_out
arrIN=(${hostname//./ })
tLen=${#arrIN[@]}
if [ $tLen -le 1 ] # Did the command work?
then # Fail
     echo "An error occured. Error: check your hostname."
     exit 1
fi

#verify LVM
#If the line starts with UUID=xyz, this means it's a physical partition.
#If the line starst with /dev/sdaX, it also means it's a physical partition.
#The indicator for LVM would be something with /dev/mapper/xyz.
cat /etc/fstab  >> verify_out

brctl show >> verify_out






#!/bin/bash

ONE_VERIFY_LOG="/var/log/megam/megamcib/one_verify.log"

kvm=`kvm-ok  | grep "KVM acceleration can be used"`

echo "$kvm" >> $ONE_VERIFY_LOG

if [ "$kvm" != "KVM acceleration can be used" ];
then # Fail
     echo "KVM isn't capable of running hw accelerate KVM virtual machines." >> $ONE_VERIFY_LOG
     exit 2
fi

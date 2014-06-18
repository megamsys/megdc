#!/bin/bash

#Verify if kvm exists
kvm=`kvm-ok  | grep "KVM acceleration can be used"`
echo "$kvm" >> verify_out

if [ "$kvm" != "KVM acceleration can be used" ]; # Did the command work?
then # Fail
     echo "An error occured. Error: KVM doesn't exist"
     exit 1
fi

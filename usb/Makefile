#!/usr/bin/make -f
#
# Cloud In a Box  Installer system makefile.
# Copyright 2012-2014 by the Megam team.
# Licensed under the terms of the Apache License V2.
#

SHELL = /bin/sh

#### Start of system configuration section. ####
SRCDIR = ../

# Temporary BUILDDIR
BUILDDIR = ./build
USB_IMAGEDIR = $(BUILDDIR)/usb_image
MOUNT_ISODIR = /mnt

TRUSTY_MINI_ISO = $(BUILDDIR)/trusty_mini.iso
TRUSTY_MEGAM_ISO = $(BUILDDIR)/trusty_megam.iso
TRUSTY_MINI_URL = "http://archive.ubuntu.com/ubuntu/dists/trusty/main/installer-amd64/current/images/netboot/mini.iso"

CENTOS_MINI_ISO = $(BUILDDIR)/centos_mini.iso
CENTOS_MEGAM_ISO = $(BUILDDIR)/centos_megam.iso
CENTOS_MINI_URL = "http://centosmirror.go4hosting.in/centos/6.5/isos/x86_64/CentOS-6.5-x86_64-minimal.iso"

TEMP_MINIISO=$(BUILDDIR)/trusty_mini.iso
MEGAM_MINIISO=$(BUILDDIR)/trusty_megam.iso

# By default, we just advertise what we can do.
.PHONY: all
all: list


# List all targets useful for direct invocation.
.PHONY: list
list:
	@echo "Useful targets:"
	@echo
	@echo "list"
	@echo "realclean"
	@echo "trusy"
	@echo "centos **work in progress.**"
	@echo

# Clean all targets.
.PHONY: realclean
realclean:
	rm -rf $(BUILDDIR)


# all_build is provided automagically, but for manual invocation
# we provide a generic build rule.
.PHONY: trusty
trusty: wget_trusty tree_iso burn_usb qemu_test
	-@if [ -e $(MEGAM_MINIISO) ]; then \
		echo "--> Successfully created trusty_megam 64bit mini distro"; \
	else \
		echo "Failure : Can't create trusty_megam mini distro." ;  exit 2; \
	fi

# all_build is provided automagically, but for manual invocation
# we provide a generic build rule.
.PHONY: centos
centos: wget_centos tree_iso burn_usb qemu_test
	-@if [ -e $(MEGAM_MINIISO) ]; then \
		echo "Successfully created centos_megam 64bit mini distro"; \
	else \
		echo "Failure : Can't create centos_megam mini distro." ;  exit 2; \
	fi


# all_build is provided automagically, but for manual invocation
# we provide a generic build rule.
.PHONY: wget_trusty
wget_trusty:
	if [  -e $(TEMP_MINIISO) ]; then \
		echo "Found trusty 64 bit mini distro"; \
	else \
		mkdir -p $(BUILDDIR); \
		echo "Downloading trusty 64 bit"; \
		wget $(TRUSTY_MINI_URL) -O $(TEMP_MINIISO); \
	fi

# all_build is provided automagically, but for manual invocation
# we provide a generic build rule.
.PHONY: wget_centos
wget_centos:
	if [ -e $(TEMP_MINIISO) ]; then \
		echo "Found centos 64 bit mini distro"; \
	else \
		mkdir -p $(BUILDDIR)
		echo "Downloading centos 64 bit ( > 200MB). This is going to take a while.."; \
		wget $(CENTOS_MINI_URL) -O $(TEMP_MINIISO); \
	fi


.PHONY:  tree_iso
tree_iso:
	if [ -e $(TEMP_MINIISO) ]; then \
		sudo mount -o loop $(TEMP_MINIISO) $(MOUNT_ISODIR); \
	else \
	  echo "mini ISO not found.";  exit 1; \
	fi;\
	sudo cp -fr $(MOUNT_ISODIR) $(USB_IMAGEDIR); \
	sudo umount $(MOUNT_ISODIR);


.PHONY: burn_usb
burn_usb:
	sudo cp preseed.cfg $(USB_IMAGEDIR)
	sudo cp txt.cfg $(USB_IMAGEDIR)
	xorriso -as mkisofs -r -J -b isolinux.bin -c boot.cat \
		-no-emul-boot -boot-load-size 4 -boot-info-table \
		-eltorito-alt-boot \
		--efi-boot boot/grub/efi.img -no-emul-boot \
		-o $(MEGAM_MINIISO) $(USB_IMAGEDIR); \
	geniso_hybrid_plus_firmware_partition $(MEGAM_MINIISO); \
	while [ -z "$$CONTINUE" ]; do \
    read -r -p "--> Burn the iso in usb[/dev/sdb]. Have you plugged it in ?   [y/N] " CONTINUE; \
	done ; \
	if [ $$CONTINUE = "y" ] || [ $$CONTINUE = "Y" ]; then \
	 (echo "Exiting."; exit 1;); \
	fi;

.PHONY:  qemu_test
qemu_test:
	qemu-system-x86_64  -cdrom $(MEGAM_MINIISO)

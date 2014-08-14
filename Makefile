#Copyright (c) 2014 Megam Systems.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
###############################################################################
# Makefile to compile cib.
# lists all the dependencies for test, prod and we can run a go build aftermath.
###############################################################################
                            

#CIBCODE_HOME = $(HOME)/code/megam/workspace/cloudinabox
dirname = $(shell pwd)
dir = $(shell ${dirname%/*/*/*/*})

A := $(shell echo $(dir))
$(info $(A))
CIBCODE_HOME = dir

export GOPATH=$(CIBCODE_HOME)

define HG_ERROR

FATAL: you need mercurial (hg) to download gulp dependencies.
       Check README.md for details
endef

define GIT_ERROR

FATAL: you need git to download gulp dependencies.
       Check README.md for details
endef

define BZR_ERROR

FATAL: you need bazaar (bzr) to download gulp dependencies.
       Check README.md for details
endef

.PHONY: all check-path get hg git bzr get-test get-prod test client

all: check-path get test

build: check-path get

# It does not support GOPATH with multiple paths.
check-path:
ifndef GOPATH
	@echo "FATAL: you must declare GOPATH environment variable, for more"
	@echo "       details, please check README.md file and/or"
	@echo "       http://golang.org/cmd/go/#GOPATH_environment_variable"
	@exit 1
endif
#ifneq ($(subst ~,$(HOME),$(GOPATH))/src/github.com/*/megam_bee, $(PWD))
#	@echo "FATAL: you must clone gulp inside your GOPATH To do so,"
#	@echo "       you can run go get github.com/megamsys/cloudinabox/..."
#	@echo "       or clone it manually to the dir $(GOPATH)/src/github.com/megamsys/cloudinabox"
#	@exit 1
#endif

clean:
	@/bin/rm -f -r $(CIBCODE_HOME)/pkg	
	@go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | tr ' ' '\n' |\
		grep '^.*\..*/.*$$' | grep -v 'github.com/megamsys/cloudinabox' |\
		sort | uniq | xargs -I{} rm -f -r $(CIBCODE_HOME)/src/{}	
	@go list -f '{{range .Imports}}{{.}} {{end}}' ./... | tr ' ' '\n' |\
		grep '^.*\..*/.*$$' | grep -v 'github.com/megamsys/cloudinabox' |\
		sort | uniq | xargs -I{} rm -f -r $(CIBCODE_HOME)/src/{} 
	@/bin/echo "Clean ...ok"

get: hg git bzr get-test get-prod

hg:
	$(if $(shell hg), , $(error $(HG_ERROR)))

git:
	$(if $(shell git), , $(error $(GIT_ERROR)))

bzr:
	$(if $(shell bzr), , $(error $(BZR_ERROR)))

get-test:
	@/bin/echo -n "Installing test dependencies... "
	@go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | tr ' ' '\n' |\
		grep '^.*\..*/.*$$' | grep -v 'github.com/megamsys/cloudinabox' |\
		sort | uniq | xargs go get -u >/tmp/.get-test 2>&1 || (cat /tmp/.get-test && exit 1)	
	@/bin/echo "ok"
	@rm -f /tmp/.get-test

get-prod:
	@/bin/echo -n "Installing production dependencies... "
	@go list -f '{{range .Imports}}{{.}} {{end}}' ./... | tr ' ' '\n' |\
		grep '^.*\..*/.*$$' | grep -v 'github.com/megamsys/cloudinabox' |\
		sort | uniq | xargs go get -u >/tmp/.get-prod 2>&1 || (cat /tmp/.get-prod && exit 1)
	@/bin/echo "ok"
	@rm -f /tmp/.get-prod

_go_test:
	@go test -i ./...
	@go test ./...

_gulpd_dry:
	@go build -o cib cib.go
	@sudo ./cib --config ./conf/cib.conf
	@rm -f cib

test: _go_test _gulpd_dry


client:
	@go build -o cib cib.go
	@echo "Done."


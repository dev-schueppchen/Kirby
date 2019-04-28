# TEMPLATE MAKEFILE FOR GO APPLICATIONS
# https://github.com/zekroTJA/makefiles
# 
# This Makefile is capable of:
#  - defautl go commands like build, run,
#    test and golint
#  - cross compiling (see make help)
#  - automatic versioning using git tag
#    and ldflags
#
# Covered by MIT License
#
# Copyright 2019 Ringo Hoffmann (zekro Development)
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

### NAMES AND LOCS ############################
APPNAME  = kirby
PACKAGE  = github.com/zekroTJA/$(APPNAME)
LDPAKAGE = static
CONFIG   = $(CURDIR)/config/private.config.yml
###############################################

### EXECUTABLES ###############################
GO     = go
DEP    = dep
GOLINT = golint
GREP   = grep
###############################################

# ---------------------------------------------

BIN    = $(CURDIR)/bin/$(APPNAME)
TAG    = $(shell git describe --tags)
COMMIT = $(shell git rev-parse HEAD)


ifneq ($(GOOS),)
	BIN := $(BIN)_$(GOOS)
endif

ifneq ($(GOARCH),)
	BIN := $(BIN)_$(GOARCH)
endif

ifneq ($(TAG),)
	BIN := $(BIN)_$(TAG)
endif

ifeq ($(OS),Windows_NT)
	BIN := $(BIN).exe
endif


PHONY = _make
_make: deps $(BIN) cleanup

PHONY += deps
deps:
	$(DEP) ensure -v

$(BIN):
	$(GO) build  \
		-v -o $@ -ldflags "\
			-X $(PACKAGE)/$(LDPAKAGE).AppVersion=$(TAG) \
			-X $(PACKAGE)/$(LDPAKAGE).AppCommit=$(COMMIT) \
			-X $(PACKAGE)/$(LDPAKAGE).Release=TRUE" \
		$(CURDIR)/cmd/$(APPNAME)

PHONY += test
test:
	$(GO) test -v -cover ./...

PHONY += lint
lint:
	$(GOLINT) ./... | $(GREP) -v vendor || true

PHONY += run
run:
	$(GO) run -v \
		$(CURDIR)/cmd/$(APPNAME) -c $(CONFIG) -loglvl 5

PHONY += cleanup
cleanup:

PHONY += help
help:
	@echo "Available targets:"
	@echo "  #        - creates binary in ./bin"
	@echo "  cleanup  - tidy up temporary stuff created by build or scripts"
	@echo "  deps     - ensure dependencies are installed"
	@echo "  lint     - run linters (golint)"
	@echo "  run      - debug run app (go run) with test config"
	@echo "  test     - run tests (go test)"
	@echo ""
	@echo "Cross Compiling:"
	@echo "  (env GOOS=linux GOARCH=arm make)"
	@echo ""
	@echo "Use different configs for run:"
	@echo "  make CONF=./myCustomConfig.yml run"
	@echo ""


.PHONY: $(PHONY)

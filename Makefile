PACKAGE      = blueprint
VERSION      = $(shell git log -n1 --pretty='%h')
BUILD_DIR    = build
RELEASE_DIR  = dist
RELEASE_FILE = $(PACKAGE)_$(VERSION)_$(shell go env GOOS)-$(shell go env GOARCH)

.PHONY: all clean clean_build clean_dist dist build install test


all: test install dist



dist: build
	mkdir -p $(RELEASE_DIR)
	mkdir -p $(BUILD_DIR)/licenses
	cp $(GOPATH)/src/github.com/pkg/browser/LICENSE $(BUILD_DIR)/licenses/pkg.browser.LICENSE
	cp $(GOROOT)/LICENSE $(BUILD_DIR)/licenses/golang.LICENSE
	cp LICENSE $(BUILD_DIR)/licenses/blueprint.LICENSE
	tar -cvzf  $(RELEASE_DIR)/$(RELEASE_FILE).tar.gz $(BUILD_DIR) --transform='s/$(BUILD_DIR)/$(RELEASE_FILE)/g'

build: clean_build
	mkdir -p $(BUILD_DIR)
	cd $(BUILD_DIR) && \
	go build github.com/urld/blueprint/cmd/blueprint && \
	go build github.com/urld/blueprint/cmd/blueprint-export

test:
	go test github.com/urld/blueprint/...


install: test
	go install github.com/urld/blueprint/cmd/blueprint
	go install github.com/urld/blueprint/cmd/blueprint-export


clean: clean_build clean_dist


clean_build:
	rm -rf $(BUILD_DIR)


clean_dist:
	rm -rf $(RELEASE_DIR)


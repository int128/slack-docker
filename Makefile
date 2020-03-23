TARGET := slack-docker
VERSION ?= latest
GITHUB_USERNAME := int128
GITHUB_REPONAME := slack-docker
LDFLAGS := -X main.version=$(VERSION)
OSARCH := linux_arm64 linux_amd64 darwin_amd64 windows_amd64

$(TARGET):
	go build -o $@ -ldflags "$(LDFLAGS)"

.PHONY: check
check:
	golangci-lint run
	go test -v ./...

.PHONY: dist
dist:
	# make the zip files for GitHub Releases
	VERSION=$(VERSION) CGO_ENABLED=0 goxzst -d dist -i "LICENSE" -o "$(TARGET)" -t "$(TARGET).rb" -osarch "$(OSARCH)" -- -ldflags "$(LDFLAGS)"
	# test the zip file
	zipinfo dist/$(TARGET)_linux_amd64.zip

.PHONY: release
release: dist
	# publish the binaries
	ghcp release -u "$(GITHUB_USERNAME)" -r "$(GITHUB_REPONAME)" -t "$(VERSION)" dist/
	# publish the Homebrew formula
	ghcp commit -u "$(GITHUB_USERNAME)" -r "homebrew-$(GITHUB_REPONAME)" -b "bump-$(VERSION)" -m "Bump the version to $(VERSION)" -C dist/ $(TARGET).rb
	ghcp pull-request -u "$(GITHUB_USERNAME)" -r "homebrew-$(GITHUB_REPONAME)" -b "bump-$(VERSION)" --title "Bump the version to $(VERSION)"

.PHONY: clean
clean:
	-rm $(TARGET)
	-rm -r dist/

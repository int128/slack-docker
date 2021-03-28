TARGET := slack-docker
TARGET_ARCHIVE := $(TARGET)_$(GOOS)_$(GOARCH).zip
TARGET_DIGEST := $(TARGET)_$(GOOS)_$(GOARCH).zip.sha256

# extract version from tag or default to latest
ifeq ($(dir $(GITHUB_REF)), refs/tags/)
  VERSION := $(notdir $(GITHUB_REF))
else
  VERSION := latest
endif

LDFLAGS := -X main.version=$(VERSION)

all: $(TARGET)

$(TARGET):
	go build -o $@ -ldflags "$(LDFLAGS)"

.PHONY: dist
dist: $(TARGET_ARCHIVE) $(TARGET_DIGEST)
$(TARGET_ARCHIVE): $(TARGET)
ifeq ($(GOOS), windows)
	powershell Compress-Archive -Path $(TARGET),LICENSE,README.md -DestinationPath $@
else
	zip $@ $(TARGET) LICENSE README.md
endif

$(TARGET_DIGEST): $(TARGET_ARCHIVE)
ifeq ($(GOOS), darwin)
	shasum -a 256 -b $(TARGET_ARCHIVE) > $@
else
	sha256sum -b $(TARGET_ARCHIVE) > $@
endif

.PHONY: dist-release
dist-release: dist
	gh release upload $(VERSION) $(TARGET_ARCHIVE) $(TARGET_DIGEST) --clobber

.PHONY: clean
clean:
	-rm $(TARGET)

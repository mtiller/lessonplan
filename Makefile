VERSION=v1.0.0
all: darwin linux windows

# Windows Targets
windows: windows-amd64

# Windows Architectures
windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o lessonplan-$(VERSION)-windows-amd64 .

# Linux Targets
linux: linux-amd64

# Linux Architectures
linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o lessonplan-$(VERSION)-linux-amd64 .

# Mac Targets
darwin: darwin-arm64

# Mac Architectures
darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o lessonplan-$(VERSION)-darwin-amd64 .

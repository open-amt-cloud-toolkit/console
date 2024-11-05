# Get version from the first argument
version=$1

docker build -t vprodemo.azurecr.io/console:v$version .

# Build for Linux
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'github.com/open-amt-cloud-toolkit/console/internal/app.Version=$version'" -trimpath -o console_linux_x64 ./cmd/app/main.go

# Build for Windows
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X 'github.com/open-amt-cloud-toolkit/console/internal/app.Version=$version'" -trimpath -o console_windows_x64.exe ./cmd/app/main.go

# Build for Mac (x64)
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X 'github.com/open-amt-cloud-toolkit/console/internal/app.Version=$version'" -trimpath -o console_mac_x64 ./cmd/app/main.go

# Build for Mac (arm64)
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -X 'github.com/open-amt-cloud-toolkit/console/internal/app.Version=$version'" -trimpath -o console_mac_arm64 ./cmd/app/main.go

# Mark the Unix system outputs as executable
chmod +x console_linux_x64
chmod +x console_mac_x64
chmod +x console_mac_arm64

# Add them to tar files respectively
tar cvfpz console_linux_x64.tar.gz console_linux_x64
tar cvfpz console_mac_x64.tar.gz console_mac_x64
tar cvfpz console_mac_arm64.tar.gz console_mac_arm64

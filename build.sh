mkdir build/
CGO_ENABLED=0 GOOS=linux go build -o build/bot.linux .
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/bot.linux-arm64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/bot.macos-arm64 .
CGO_ENABLED=0 GOOS=windows go build -o build/bot.windows.exe .
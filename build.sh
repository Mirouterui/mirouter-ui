echo "Building MirouterUI..."
echo "Deleting old build..."
rm mirouterui_win_amd64.exe
rm mirouterui_win_386.exe
rm mirouterui_win_arm64.exe
rm mirouterui_linux_amd64
rm mirouterui_linux_arm64
rm mirouterui_linux_mipsle

echo "Building win_amd64"
GOOS=windows GOARCH=amd64 go build main.go
mv main.exe mirouterui_win_amd64.exe
upx --best -o mirouterui_win_amd64_upx.exe mirouterui_win_amd64.exe
echo "Building win_386"
GOOS=windows GOARCH=386 go build main.go
mv main.exe mirouterui_win_386.exe
upx --best -o mirouterui_win_386_upx.exe mirouterui_win_386.exe
echo "Building win_arm64"
GOOS=windows GOARCH=arm64 go build main.go
mv main.exe mirouterui_win_arm64.exe
# rm main.exe

echo "Building linux_amd64"
GOOS=linux GOARCH=amd64 go build main.go
mv main main_linux_amd64
upx --best -o mirouterui_linux_amd64_upx main_linux_amd64
echo "Building linux_arm64"
GOOS=linux GOARCH=arm64 go build main.go
mv main main_linux_arm64
upx --best -o mirouterui_linux_arm64_upx main_linux_arm64
echo "Building linux_mipsle"
GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build main.go
mv main main_linux_mipsle
upx --best -o mirouterui_linux_mipsle_upx main_linux_mipsle
rm main

echo "Building darwin_amd64"
GOOS=darwin GOARCH=amd64 go build main.go
mv main main_darwin_amd64
upx --best -o mirouterui_darwin_amd64_upx main_darwin_amd64
echo "Building darwin_arm64"
GOOS=darwin GOARCH=arm64 go build main.go
mv main main_darwin_arm64
upx --best -o mirouterui_darwin_arm64_upx main_darwin_arm64
rm main


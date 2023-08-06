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
#mv main.exe mirouterui_win_amd64.exe
upx --best -o mirouterui_win_amd64.exe main.exe
echo "Building win_386"
GOOS=windows GOARCH=386 go build main.go
# mv main.exe mirouterui_win_386.exe
upx --best -o mirouterui_win_386.exe main.exe
echo "Building win_arm64"
GOOS=windows GOARCH=arm64 go build main.go
mv main.exe mirouterui_win_arm64.exe
# rm main.exe

echo "Building linux_amd64"
GOOS=linux GOARCH=amd64 go build main.go
# mv main main_linux_amd64
upx --best -o mirouterui_linux_amd64 main
echo "Building linux_arm64"
GOOS=linux GOARCH=arm64 go build main.go
# mv main main_linux_arm64
upx --best -o mirouterui_linux_arm64 main
echo "Building linux_mipsle"
GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build main.go
# mv main main_linux_mipsle
upx --best -o mirouterui_linux_mipsle main
rm main
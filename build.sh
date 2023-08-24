VERSION="1.1.0"
echo "Building MirouterUI $VERSION..."
echo "Deleting old build..."
rm mirouterui_win_amd64.exe
rm mirouterui_win_386.exe
rm mirouterui_win_arm64.exe
rm mirouterui_linux_amd64
rm mirouterui_linux_arm64
rm mirouterui_linux_mipsle
rm mirouterui_linux_mips
rm mirouterui_linux_armv5
rm mirouterui_linux_armv6
rm mirouterui_linux_armv7
rm mirouterui_darwin_amd64
rm mirouterui_darwin_arm64

echo "Building win_amd64"
GOOS=windows GOARCH=amd64 go build main.go
mv main.exe mirouterui_win_amd64_noupx_$VERSION.exe
upx --best -o mirouterui_win_amd64_$VERSION.exe mirouterui_win_amd64_noupx_$VERSION.exe

echo "Building win_386"
GOOS=windows GOARCH=386 go build main.go
mv main.exe mirouterui_win_386_noupx_$VERSION.exe
upx --best -o mirouterui_win_386_$VERSION.exe mirouterui_win_386_noupx_$VERSION.exe

echo "Building win_arm64"
GOOS=windows GOARCH=arm64 go build main.go
mv main.exe mirouterui_win_arm64_noupx_$VERSION.exe
# rm main.exe

echo "Building linux_amd64"
GOOS=linux GOARCH=amd64 go build main.go
mv main mirouterui_linux_amd64_noupx_$VERSION
upx --best -o mirouterui_linux_amd64_$VERSION mirouterui_linux_amd64_noupx_$VERSION

echo "Building linux_arm64"
GOOS=linux GOARCH=arm64 go build main.go
mv main mirouterui_linux_arm64_noupx_$VERSION
upx --best -o mirouterui_linux_arm64_$VERSION mirouterui_linux_arm64_noupx_$VERSION

echo "Building linux_mipsle"
GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build main.go
mv main mirouterui_linux_mipsle_noupx_$VERSION
upx --best -o mirouterui_linux_mipsle_$VERSION mirouterui_linux_mipsle_noupx_$VERSION

echo "Building linux_mips"
GOOS=linux GOARCH=mips GOMIPS=softfloat go build main.go
mv main mirouterui_linux_mips_noupx_$VERSION
upx --best -o mirouterui_linux_mips_$VERSION mirouterui_linux_mips_noupx_$VERSION

echo "Building linux_armv5"
GOOS=linux GOARCH=arm GOARM=5 GOMIPS=softfloat go build main.go
mv main mirouterui_linux_armv5_noupx_$VERSION
upx --best -o mirouterui_linux_armv5_$VERSION mirouterui_linux_armv5_noupx_$VERSION

echo "Building linux_armv6"
GOOS=linux GOARCH=arm GOARM=6 GOMIPS=softfloat go build main.go
mv main mirouterui_linux_armv6_noupx_$VERSION
upx --best -o mirouterui_linux_armv6_$VERSION mirouterui_linux_armv6_noupx_$VERSION

echo "Building linux_armv7"
GOOS=linux GOARCH=arm GOARM=7 GOMIPS=softfloat go build main.go
mv main mirouterui_linux_armv7_noupx_$VERSION
upx --best -o mirouterui_linux_armv7_$VERSION mirouterui_linux_armv7_noupx_$VERSION

echo "Building darwin_amd64"
GOOS=darwin GOARCH=amd64 go build main.go
mv main mirouterui_darwin_amd64_noupx_$VERSION
upx --best -o mirouterui_darwin_amd64_$VERSION mirouterui_darwin_amd64_noupx_$VERSION

echo "Building darwin_arm64"
GOOS=darwin GOARCH=arm64 go build main.go
mv main mirouterui_darwin_arm64_noupx_$VERSION
upx --best -o mirouterui_darwin_arm64_$VERSION mirouterui_darwin_arm64_noupx_$VERSION
# rm main


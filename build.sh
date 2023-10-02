if [ -z "$1" ]
then
  read -p "版本号: " VERSION
else
  VERSION=$1
fi
echo "版本号是：$VERSION"
echo "Building MirouterUI $VERSION..."
OUTPUT_DIR="./build"
echo "编译生成目录：$OUTPUT_DIR"
mkdir -p $OUTPUT_DIR

if [ -z "$2" ] || [ "$2" = "win" ]
then
    echo "Building win_amd64"
    GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_win_amd64_noupx_$VERSION.exe main.go
    upx --best -o $OUTPUT_DIR/mirouterui_win_amd64_$VERSION.exe $OUTPUT_DIR/mirouterui_win_amd64_noupx_$VERSION.exe

    echo "Building win_386"
    GOOS=windows GOARCH=386 go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_win_386_noupx_$VERSION.exe main.go
    upx --best -o $OUTPUT_DIR/mirouterui_win_386_$VERSION.exe $OUTPUT_DIR/mirouterui_win_386_noupx_$VERSION.exe

    echo "Building win_arm64"
    GOOS=windows GOARCH=arm64 go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_win_arm64_noupx_$VERSION.exe main.go
    # rm main.exe
fi

if [ -z "$2" ] || [ "$2" = "linux" ]
then
  echo "Building linux_amd64"
  GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_linux_amd64_noupx_$VERSION main.go
  upx --best -o $OUTPUT_DIR/mirouterui_linux_amd64_$VERSION $OUTPUT_DIR/mirouterui_linux_amd64_noupx_$VERSION

  echo "Building linux_arm64"
  GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_linux_arm64_noupx_$VERSION main.go
  upx --best -o $OUTPUT_DIR/mirouterui_linux_arm64_$VERSION $OUTPUT_DIR/mirouterui_linux_arm64_noupx_$VERSION

  echo "Building linux_mipsle"
  GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_linux_mipsle_noupx_$VERSION main.go
  upx --best -o $OUTPUT_DIR/mirouterui_linux_mipsle_$VERSION $OUTPUT_DIR/mirouterui_linux_mipsle_noupx_$VERSION

  echo "Building linux_mips"
  GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_linux_mips_noupx_$VERSION main.go
  upx --best -o $OUTPUT_DIR/mirouterui_linux_mips_$VERSION $OUTPUT_DIR/mirouterui_linux_mips_noupx_$VERSION

  echo "Building linux_s390x"
  GOOS=linux GOARCH=s390x go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_linux_s390x_noupx_$VERSION main.go
  # upx --best -o $OUTPUT_DIR/mirouterui_linux_s390x_$VERSION $OUTPUT_DIR/mirouterui_linux_s390x_noupx_$VERSION
  for version in 5 6 7
  do
      echo "Building linux_armv$version"
      GOOS=linux GOARCH=arm GOARM=$version GOMIPS=softfloat go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_linux_armv${version}_noupx_$VERSION main.go
      upx --best -o $OUTPUT_DIR/mirouterui_linux_armv${version}_$VERSION $OUTPUT_DIR/mirouterui_linux_armv${version}_noupx_$VERSION
  done
fi


if [ -z "$2" ] || [ "$2" = "darwin" ]
then
# Building darwin_amd64
echo "Building darwin_amd64"
GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_darwin_amd64_noupx_$VERSION main.go
upx --best -o $OUTPUT_DIR/mirouterui_darwin_amd64_$VERSION $OUTPUT_DIR/mirouterui_darwin_amd64_noupx_$VERSION

# Building darwin_arm64
echo "Building darwin_arm64"
GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.Version=$VERSION'" -o $OUTPUT_DIR/mirouterui_darwin_arm64_noupx_$VERSION main.go
upx --best -o $OUTPUT_DIR/mirouterui_darwin_arm64_$VERSION $OUTPUT_DIR/mirouterui_darwin_arm64_noupx_$VERSION
fi
# pixel-experiments/pong/pong-with-sound

http://cs.au.dk/~dsound/DigitalAudio.dir/Greenfoot/Pong.dir/Pong.html

## Cross compile from macOS to Windows

`CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build pong-with-sound.go`

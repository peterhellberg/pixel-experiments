go get gopkg.in/slimsag/file2go.v1

lame -b 32 --resample 32 -a beeep.wav beeep.mp3
file2go.v1 -f -i beeep.mp3 -o beeep.go -package main -var beeep

lame -b 32 --resample 32 -a peeeeeep.wav peeeeeep.mp3
file2go.v1 -f -i peeeeeep.mp3 -o peeeeeep.go -package main -var peeeeeep

lame -b 32 --resample 32 -a plop.wav plop.mp3
file2go.v1 -f -i plop.mp3 -o plop.go -package main -var plop

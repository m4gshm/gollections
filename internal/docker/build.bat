
git clone https://github.com/golang/go.git
cd go
git checkout release-branch.go1.18
del /Q bin
del /Q pkg
cd src
SET GOOS=linux
SET GOARCH=amd64
make.bat

docker build -t m4gshm/golang:1.18rc1 .
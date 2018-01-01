# lxdlet
Lxd for Kubernetes. This project is under development.

## Setup your development environment
This project uses go 1.9.

Install go land on Ubuntu Xenial:
```
$ sudo add-apt-repository ppa:gophers/archive
$ sudo apt update
$ sudo apt-get install golang-1.9-go
```

Prepare your bash environment:
```
$ echo "export GOROOT=/usr/lib/go-1.9/" >> ~/.bashrc
$ mkdir $HOME/go-workspace
$ echo "export GOPATH=\$HOME/go-workspace/" >> ~/.bashrc
$ echo "export PATH=/usr/lib/go-1.9/bin/:\$GOPATH:\$PATH" >> ~/.bashrc
$ source ~/.bashrc
```

Get the project:
```
$ go get github.com/ktsakalozos/lxdlet
```

Build lxdlet:
```
go build -o lxdlet  $GOPATH/src/github.com/ktsakalozos/lxdlet/lxdlet/cmd/server/main.go
```

# Resources
Golang on Ubuntu: https://github.com/golang/go/wiki/Ubuntu

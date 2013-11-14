
build: prebuild
	go build

prebuild:
	go get github.com/jteeuwen/go-bindata
	go-bindata -toc -prefix "$$PWD/plugins" -out "happy.js.go" plugins/happy.js

export PATH := $(GOPATH)/bin:${PATH}

all: cross build

build:
	go build -o bin/b2 ./cli

cross:
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_darwin_amd64.tar.gz ./bin/b2
	env GOOS=freebsd GOARCH=386 go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_freebsd_386.tar.gz ./bin/b2
	env GOOS=freebsd GOARCH=amd64 go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_freebsd_amd64.tar.gz ./bin/b2
	env GOOS=linux GOARCH=386 go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_linux_386.tar.gz ./bin/b2
	env GOOS=linux GOARCH=amd64 go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_linux_amd64.tar.gz ./bin/b2
	env GOOS=linux GOARCH=arm go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_linux_arm.tar.gz ./bin/b2
	env GOOS=linux GOARCH=arm64 go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_linux_arm64.tar.gz ./bin/b2
	env GOOS=windows GOARCH=386 go build -o ./bin/b2.exe ./cli
	tar zcf ./bin/b2_${B2_VERSION}_windows_386.tar.gz ./bin/b2
	env GOOS=windows GOARCH=amd64 go build -o ./bin/b2.exe ./cli
	tar zcf ./bin/b2_${B2_VERSION}_windows_amd64.tar.gz ./bin/b2
	env GOOS=linux GOARCH=mips64 go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_linux_mips64.tar.gz ./bin/b2
	env GOOS=linux GOARCH=mips64le go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_linux_mips64le.tar.gz ./bin/b2
	env GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_linux_mips.tar.gz ./bin/b2
	env GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -o ./bin/b2 ./cli
	tar zcf ./bin/b2_${B2_VERSION}_linux_mipsle.tar.gz ./bin/b2

test:
	go test

clean:
	rm -fr bin

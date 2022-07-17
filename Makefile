
build:
	go build cmd/golox/golox.go

test-suite: build
	pushd ../craftinginterpreters/ && dart tool/bin/test.dart jlox --interpreter ../crafting-interpreters-go/golox && popd

clean:
	rm golox

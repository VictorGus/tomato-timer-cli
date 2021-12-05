run:
	cd cmd/tomato-timer && go run main.go
build:
	cd cmd/tomato-timer && go build 
pack:
	rm -rf target/tomato-timer && mkdir target/tomato-timer && make build && cp cmd/tomato-timer/tomato-timer target/tomato-timer/tomato-timer && cp -rf resources target/tomato-timer/resources
build:
	go build -o ./build/godepuml cmd/godepuml/main.go

clean:
	rm -rf ./build
	rm *.puml
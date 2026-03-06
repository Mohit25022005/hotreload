build:
	go build -o hotreload.exe ./cmd/hotreload

run:
	.\hotreload.exe --root ./testserver --build "go build -o ./bin/server.exe ./testserver" --exec ./bin/server.exe

test:
	go test ./...
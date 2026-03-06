build:
	go build -o hotreload ./cmd/hotreload

demo:
	go build -o ./bin/server ./testserver
	./hotreload --root ./testserver \
	--build "go build -o ./bin/server ./testserver" \
	--exec "./bin/server"
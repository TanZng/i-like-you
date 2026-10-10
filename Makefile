.PHONY: build serve clean

build:
	GOOS=js GOARCH=wasm go build -o web/game.wasm .

serve: build
	cd web && python3 -m http.server 8080

clean:
	rm -f web/game.wasm

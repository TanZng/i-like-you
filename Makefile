.PHONY: build serve clean

build:
	GOOS=js GOARCH=wasm go build -o web/game.wasm .
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/wasm_exec.js

serve: build
	cd web && python3 -m http.server 8080

clean:
	rm -f web/game.wasm web/wasm_exec.js

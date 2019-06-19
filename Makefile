main: clean wasm_exec
	tinygo build -o ./html/wasm.wasm -target wasm -no-debug ./main/main.go
	cp ./main/index.html ./html/

wasm_exec:
	cp ../../../targets/wasm_exec.js ./html/

clean:
	rm -rf ./html
	mkdir ./html


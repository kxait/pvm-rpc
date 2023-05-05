build:
	which gcc > /dev/null
	which pvm > /dev/null
	which pvmgetarch > /dev/null
	which go > /dev/null
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/pvm-rpc .

copy:
	which pvm > /dev/null
	which pvmgetarch > /dev/null
	test -n "$(PVM_PATH)" # $$PVM_PATH
	cp bin/linux-amd64/pvm-rpc ${PVM_PATH}/

clean:
	rm -rf bin

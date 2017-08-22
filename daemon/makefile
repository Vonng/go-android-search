run:
	go run android.go

start:
	nohup ./android 1>android.log 2>&1 &

stop:
	ps aux | grep "\./android" | grep -v grep | awk '{print $2}' | xargs kill

status:
	ps aux | grep "\./android"

log:
	tail -f android.log

build:
	go build -o android android.go
	chmod a+x android

install:
	go build -o android android.go
	chmod a+x android
	mv android $(GOPATH)/bin/android

upload:
	GOOS=linux GOARCH=amd64 go build -o android android.go
	chmod a+x android
	oss put android && rm android

download:
	oss get android
	chmod a+x android

linux:
	GOOS=linux GOARCH=amd64 go build -o android android.go
	chmod a+x android

mac:
	GOOS=darwin GOARCH=amd64 go build -o android android.go
	chmod a+x android

createdb:
	psql -c "CREATE USER meta CREATEDB;"
	psql -c "ALTER USER meta with PASSWORD 'meta';"
	psql postgres meta -c "CREATE DATABASE meta;"

setup:
	psql meta meta < android.ddl

clean:
	rm -rf android android.log

.PHONY: run clean start stop log build install upload download linux mac createdb setup


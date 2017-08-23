# android daemon

android daemon that take ID as input, produce android Application as output.


## Install

```bash
go get github.com/Vonng/go-android-search
cd ${GOPATH}/src/github.com/Vonng/go-android-search/daemon

# Setup database environment. assume you have an available local pg
# It will create a user `meta` with owns a database named `meta`
make createdb

# It will create table `android` and `android_queue` in database `meta`
make setup

# Build binary
make build

# Install: mv binary to your $GOPATH
make install
```

now everything is prepared for running the daemon

### Usage

some frequently used bash command can be accessed from makefile

```bash
# Start the daemon. don't forget build before start
make start

# show daemon status
make status

# Stop the daemon
make stop

# See log
make log

# Using `go run android.go`
make
```


### Assign Task

INSERT into `android_queue`. `android` will take task from queue table and put result into table `android`.
 
 task format is `TypeLetter + ID`, where `TypeLetter` could be:
 
 * `!`: stand for package name
 * `#`: stand for keyword.  program will search and fetch new found app.
 * no leading letter will use bundleID by default. (for stupid client...)


And daemon binary can handle iTunesID, BundleID, Keywords directly by:

```bash
android id com.tencent.xin
android key yourKeyword
```

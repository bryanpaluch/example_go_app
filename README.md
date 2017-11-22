# Example Golang App

```

// create the required directory and 
$ mkdir -p $GOPATH/github.com/bryanpaluch                                                                           
$ cd $GOPATH/github.com/bryanpaluch                                                                                 
// clone this repo
$ git@github.com:bryanpaluch/example_go_app.git
$ cd example_go_app
// pulls down dependencies using dep, runs tests in all packages and builds both linux and mac versions
$ make

// builds linux binary and builds docker image
$ make docker

// prepares dev envo for docker-compose volumes
$ make prepare

// spins up docker mysql and example app

$ docker-compose up

```

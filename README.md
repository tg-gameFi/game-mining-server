# TG Mini Wallet Server

## Environment

1. Golang V1.22 or above
2. Mysql V8 or above
3. Redis V6 or above

## TestLocal

0. Init Dependencies:  `make docker-up`
1. Build exec: `make build`
2. Start and Run: `make serve`

## Build & Run

Before build, assume you have been configured Golang environment or have Docker installed, Mysql and Redis are also
necessary to be installed. make sure Mysql and Redis are ready and then set mysql config and redis config in `config.{env}.yaml`

0. Init mysql & redis locally: (remove server and web) `docker-compose up -d`
1. Download dependencies: Clone and run command `go mod tidy` in project root dir
2. Build executable file: Run command `rm entry && go build -o entry .` in project root dir
3. Start server in local: Run command `./entry dev` in project root dir

Supported env args: `dev`, `test`, `prod`, for different env, it will use
different config files: `config.dev.yaml`, `config.test.yaml`, `config.prod.yaml`, please make sure you have specified
config file in project root dir

**If you want to move `entry` outside of project root dir to run, you have to copy `config/` dir and `config.{env}.yaml`
with it**

The server needs 2 args to start: `env` and `botToken`, for local develop you may not need `botToken`, just
run `./entry dev` to start a local server without Telegram bot functions

## API Docs

**API Docs only host in `dev` and `test` env**

Execute command `swag init` in project root dir to generate newly swagger api docs, it will create `swagger.json`
and `swagger.yaml` files in dir `/docs`

Run server in `dev` or `test` env, then visit `http://127.0.0.1:8080/swagger/index.html` you can see api docs page



# prooxy
>db prooxy app implements golang

## go version
>1.6

## TODO
> config

## db
```sql
CREATE TABLE redirect (
	before_id varchar PRIMARY KEY,
	after_id varchar NOT NULL
);

INSERT INTO public.redirect VALUES ('id1', 'id2');
```

## Installation
```bash
go get github.com/gin-gonic/gin
go get github.com/robfig/config
go get github.com/lib/pq
go get github.com/fvbock/endless
```

## run for dev
```bash
go run proxy.go
```

## build
```bash
go build
```

## Cross Complie
```bash
GOOS=linux GOARCH=amd64 go build proxy.go
```

## run for release mode
```bash
export GIN_MODE=release
./proxy > /var/log/proxy.log
```

@see release/start_proxy.sh

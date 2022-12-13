cli:
	go build -mod vendor -o bin/b64e cmd/b64e/main.go
	go build -mod vendor -o bin/b64d cmd/b64d/main.go
	go build -mod vendor -o bin/md5e cmd/md5e/main.go
	go build -mod vendor -o bin/urlencode cmd/urlencode/main.go
	go build -mod vendor -o bin/now cmd/now/main.go
	go build -mod vendor -o bin/when cmd/when/main.go

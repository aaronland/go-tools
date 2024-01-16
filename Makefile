GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/b64e cmd/b64e/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/b64d cmd/b64d/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/md5e cmd/md5e/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/urlencode cmd/urlencode/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/now cmd/now/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/when cmd/when/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/qrd cmd/qrd/main.go

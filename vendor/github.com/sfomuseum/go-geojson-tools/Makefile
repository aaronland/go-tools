GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/bounds cmd/bounds/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/bbox2feature cmd/bbox2feature/main.go

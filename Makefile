GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	@make build-tool TOOL=b64e
	@make build-tool TOOL=b64d
	@make build-tool TOOL=md5e
	@make build-tool TOOL=urlescape
	@make build-tool TOOL=urlunescape
	@make build-tool TOOL=now
	@make build-tool TOOL=when
	@make build-tool TOOL=ts
	@make build-tool TOOL=qrd
	@make build-tool TOOL=qre
	@make build-tool TOOL=jp
	@make build-tool TOOL=jf
	@make build-tool TOOL=jv
	@make build-tool TOOL=gh2bb
	@make build-tool TOOL=pt2gh
	@make build-tool TOOL=bb2f
	@make build-tool TOOL=fc2d
	@make build-tool TOOL=pts2ls

build-tool:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/$(TOOL) cmd/$(TOOL)/main.go

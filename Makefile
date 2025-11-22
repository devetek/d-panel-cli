install:
	@export GOPRIVATE="github.com/devetek/*" && go mod tidy

build:
	./scripts/build.sh
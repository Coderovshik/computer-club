.PHONY: build
build:
	@go build -o bin/out ./*.go

.PHONY: run
run: build
	@./bin/out

.PHONY: clean
clean:
	@rm -rf bin

.PHONY: image
image:
	@docker build -t compclub .
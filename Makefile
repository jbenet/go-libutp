all: build

deps:
	go get ./...

build:
	go build

install: build
	go install

watch:
	-make install
	@echo "[watching *.go for recompilation]"
	# for portability, use watchmedo -- pip install watchmedo
	@watchmedo shell-command --patterns="*.go;" --recursive \
		--command='make build' .

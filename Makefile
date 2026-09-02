.PHONY: gofmt
gofmt:
	gofumpt -l -w .
	golines -w --max-len=120 .

GO ?= go
TESTPKG := testextpoints
TESTTARGET := ./...

test:
	$(GO) run main.go template.go $(TESTPKG); $(GO) test -v $(TESTTARGET)

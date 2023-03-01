SHELL := /bin/bash
# placeholder becasuse some packages have init functions which check TCE_PSM env var.
psm=product.subsystem.module

export TCE_PSM=$(psm)
export GO111MODULE=on

format:
	find . -name '*.go' | grep -Ev 'vendor|thrift_gen' | xargs goimports -w

unittest:
	go run ./testrunner -tags=unit
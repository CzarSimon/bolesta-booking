run-local:
	sh run-local.sh

test:
	go test ./...

full-test:
	go clean -testcache
	go test ./...
	go vet ./...
	trivy fs .

run-local:
	sh run-local.sh

test:
	go test ./...

full-test:
	go test ./...
	go vet ./...
	trivy fs .

TEST_CMD := ./shortenertestbeta -test.v -source-path=. -binary-path=./cmd/shortener/shortener -server-port=8087 -database-dsn='postgres://app:pass@localhost:5432/short?sslmode=disable'


test:
	go test ./internal/... -v

generate-test-mocks:
	mockgen -source=internal/contracts.go -destination=mocks/postgres_mock.go -package=mocks Storage

vet:
	go vet -vettool=statictest  ./...

test-iteration:
	go build -o cmd/shortener/shortener cmd/shortener/*.go
	$(TEST_CMD) -test.run=^TestIteration1$$
	$(TEST_CMD) -test.run=^TestIteration2$$
	$(TEST_CMD) -test.run=^TestIteration3$$
	$(TEST_CMD) -test.run=^TestIteration4$$
	$(TEST_CMD) -test.run=^TestIteration5$$
	$(TEST_CMD) -test.run=^TestIteration6$$
	$(TEST_CMD) -test.run=^TestIteration7$$
	$(TEST_CMD) -test.run=^TestIteration8$$
	$(TEST_CMD) -test.run=^TestIteration9$$ -file-storage-path=.vscode/tmp.json
	$(TEST_CMD) -test.run=^TestIteration10$$ -file-storage-path=.vscode/tmp.json
	$(TEST_CMD) -test.run=^TestIteration11$$ -file-storage-path=.vscode/tmp.json
	$(TEST_CMD) -test.run=^TestIteration12$$ -file-storage-path=.vscode/tmp.json
	$(TEST_CMD) -test.run=^TestIteration13$$ -file-storage-path=.vscode/tmp.json
	$(TEST_CMD) -test.run=^TestIteration14$$ -file-storage-path=.vscode/tmp.json

up:
	docker compose up -d

down:
	docker compose down

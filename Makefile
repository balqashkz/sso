
sso-local:
	@echo "Запускаем..."
	go run ./cmd/sso/main.go --config=./config/local.yaml

sso-tests:
	@echo "Запускаем..."
	go test ./tests/auth_register_login_test.go

sso-migrate:
	@echo "Запускаем миграцию..."
	go run ./cmd/migrator/main.go \
 		--storage-path=./storage/sso.db \
 		--migrations-path=./migrations

sso-migrate-tests:
	@echo "Запускаем миграцию..."
	go run ./cmd/migrator/main.go \
        --storage-path=./storage/sso.db \
        --migrations-path=./tests/migrations \
        --migrations-table=migrations_tests
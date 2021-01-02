.PHONY: migrate-up migrate-down generate

generate:
	go generate ./...

migrate-up:
	echo "see <https://github.com/golang-migrate/migrate/issues/59>"
	#migrate -verbose -source file://scripts -database sqlite3://testdata/accoutinformation.sqlite up

migrate-down:
	echo "see <https://github.com/golang-migrate/migrate/issues/59>"
	#migrate -verbose -source file://scripts -database sqlite3://testdata/accoutinformation.sqlite up

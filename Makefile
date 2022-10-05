DB_URL = postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres14 --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:14-alpine
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
	
dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migrations -database "$(DB_URL) -verbose up 1

migratedown:
	migrate -path db/migrations -database "$(DB_URL) -verbose down

migratedown1:
	migrate -path db/migrations -database "$(DB_URL) -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./... -coverprofile cover.out

testcoverage:
	go tool cover -func cover.out | grep total | awk '{print $3}'

server:
	go run main.go
	
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/muditshukla3/simplebank/db/sqlc Store

.PHONY: postgres dropdb createdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock testcoverage
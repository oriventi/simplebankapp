initnetwork:
	docker network create bank-network

initdb:
	docker run -d --name postgres_container --network bank-network -p3808:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:alpine3.17
startdb:
	docker start postgres_container
stopdb:
	docker stop postgres_container
createdb:
	docker exec -it postgres_container createdb simple_bank
dropdb:
	docker exec -it postgres_container dropdb simple_bank

test:
	go test -v -cover ./...
server:
	go run main.go
initdockerserver:
	make stopdb
	make startdb
	docker run --name simplebank --network bank-network -p8080:8080 -e DB_SOURCE="postgresql://root:secret@postgres_container/simple_bank?sslmode=disable" simplebank:latest
startdockerserver:
	docker start simplebank
stopdockerserver:
	docker stop simplebank

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema
migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:3808/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:3808/simple_bank?sslmode=disable" -verbose down
migrateup1:
	migrate -path db/migration -database "postgres://root:secret@localhost:3808/simple_bank?sslmode=disable" -verbose up 1
migratedown1:
	migrate -path db/migration -database "postgres://root:secret@localhost:3808/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v "C:/Users/adria/Desktop/Programmieren/golang/simplebankapp:/src" -w /src kjconroy/sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/oriventi/simplebank/db/sqlc Store

dockerimg:
	docker build -t simplebank:latest .

.PHONY: startdockerserver stopdockerserver initdockerserver dockerimg createdb dropdb startdb stopdb initdb migrateup migratedown migrateup1 migratedown1 sqlc migratecreate test server mock

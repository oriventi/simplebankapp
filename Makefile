initdb:
	docker run -d --name postgres_container -p3808:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:alpine3.17
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

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema
migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:3808/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:3808/simple_bank?sslmode=disable" -verbose down
sqlc:
	docker run --rm -v "C:/Users/adria/Desktop/Programmieren/golang/simplebankapp:/src" -w /src kjconroy/sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/oriventi/simplebank/db/sqlc Store

.PHONY: createdb dropdb startdb stopdb initdb migrateup migratedown sqlc migratecreate test server mock

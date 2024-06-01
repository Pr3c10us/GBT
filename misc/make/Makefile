include ../../.dev.env

define db_up
	    migrate -path ../../stack/migrations/postgres -database "postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_HOST}:${PG_DB_PORT}/${PG_DB_NAME}?sslmode=${PG_SSL_MODE}"  up
endef

define db_down
	    migrate -path ../../stack/migrations/postgres -database "postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_HOST}:${PG_DB_PORT}/${PG_DB_NAME}?sslmode=${PG_SSL_MODE}"  down
endef

define db_force
	    migrate -database "postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_HOST}:${PG_DB_PORT}/${PG_DB_NAME}?sslmode=${PG_SSL_MODE}" -path ../../stack/migrations/postgres force $(version)
endef

define db_postgresql_create
	migrate create -ext sql -dir ../../stack/migrations/postgres -seq $(migration)
endef

pg_db_force:
	$(call db_force)

pg_db_down:
	$(call db_down)

pg_db_up:
	$(call db_up)

pg_db_create:
	$(call db_postgresql_create)

generate_datastore: pg_db_up
	sqlc generate -f infrastructures/migrations/postgresqlsqlc.yml

dev: pg_db_up
	@air

build_http: pg_db_up
	@go build -o bin/http/main ./cmd/http/main.go

run_http: build_http
	@bin/http/main
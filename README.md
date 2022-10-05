
### Commands

make postgres - To start the postgres container

make createdb - To create the database

### Migrations

Install golang-migrate tool to run migrations

```
brew install golang-migrate
```

Once installation is complete run the following command to create migration

```
migrate create -ext sql -dir db/migrations -seq init_schema
```

Above command will create init_schema.down.sql and init_schema.up.sql.

Running the migrations

```
migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disabled" -verbose up
```

## Authorization Rules

API Create Account - A logged-in user can only create an account for him/herself
API Get Account - A logged-in user can only get account that he/she owns.
API List Account - A logged-in user can only list accounts that belong to him/herself.
Api Transfer Money - A logged-in user can only send money from his/her own account.
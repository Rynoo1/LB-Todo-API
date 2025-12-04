# LB-Todo-API

A simple TODO API built in Go, using PostgreSQL for storage.

## Prerequisites

- Go (e.g. 1.20+)  
- PostgreSQL (local or Docker)  
- `git`  

## Setup — Local (PostgreSQL)

1. Clone the repo  
   ```bash
   git clone https://github.com/Rynoo1/LB-Todo-API.git
   cd LB-Todo-API
   
2. Configure database credentials
- Copy the `.env` file or create a new `.env` in project root
- Update the environment variables (e.g. `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`) to point to your local Postgres setup

3. Create the database
```bash
psql -U <your-db-user> -h <host> -p <port> -c "CREATE DATABASE <your-db-name>;"
```
5. Start the server
```bash
go run main.go
```

The server will start (e.g. on `localhost:8080`)

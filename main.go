package main

import (
	"log"

	"github.com/Rynoo1/LB-Todo-API/config"
	"github.com/Rynoo1/LB-Todo-API/migrate"
	"github.com/gofiber/fiber/v2"
)

func main() {

	// init db
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	// get sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("error getting sql.DB : %v", err)
	}

	// ping db
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}
	log.Println("database connection OK")

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatalf("error closing database: %v", err)
		}
	}()

	if err := migrate.RunMigrations(db); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	println("migrations run successfully!")

	app := fiber.New()

	log.Fatal(app.Listen(":8080"))

}

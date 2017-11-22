package main

import (
	"fmt"
	"github.com/bryanpaluch/example_go_app/controller"
	"github.com/bryanpaluch/example_go_app/db"
	"log"
	"os"
	"time"
)

var (
	version, commit, branch string
)

func GetEnv(key string, def string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return def
}

func main() {
	start := time.Now()
	log.Printf("starting example service, version=%s, commit=%s, branch=%s", version, commit, branch)

	exampleDb, err := db.NewExampleDB(fmt.Sprintf("%s:%s@(%s:3306)/mysql?parseTime=true",
		GetEnv("DB_USER", "root"),
		GetEnv("DB_PASS", "secret"),
		GetEnv("DB_HOST", "localhost")))
	if err != nil {
		log.Fatal("failed to create new db", err)
	}

	migrationPath := GetEnv("MIGRATION_PATH", "/usr/data/migrations")
	err = exampleDb.ConnectAndMigrate(migrationPath)
	if err != nil {
		log.Fatal("failed to connect or migrate", err)
	}
	router, err := controller.NewRouter(exampleDb)

	go router.Start()

	log.Printf("started server in %s", time.Since(start).String())
	select {}
}

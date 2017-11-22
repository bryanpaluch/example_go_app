package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"os"
	"testing"
)

var db *sql.DB
var dsn string

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		dsn = fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp"))
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestMigrate(t *testing.T) {
	db, _ := NewExampleDB(dsn)
	err := db.ConnectAndMigrate("./migrations")
	if err != nil {
		t.Error("failed to connect and migrate", err)
		t.Fail()
	}
}

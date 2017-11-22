package db

import (
	"context"
	"database/sql"
	"github.com/cenk/backoff"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
)

func migrateDB(db *sql.DB, directory string) error {
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+directory,
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

type DB interface {
	GetPersonByID(ctx context.Context, id int) (*Person, error)
}

type ExampleDB struct {
	dsn string
	*sql.DB
}

func NewExampleDB(dsn string) (*ExampleDB, error) {
	return &ExampleDB{dsn, nil}, nil
}

func (edb *ExampleDB) ConnectAndMigrate(directory string) error {
	err := backoff.Retry(func() error {
		conn, err := sql.Open("mysql", edb.dsn)
		if err != nil {
			return err
		}
		err = conn.Ping()
		if err != nil {
			return err
		}
		edb.DB = conn
		return nil
	}, backoff.WithMaxTries(backoff.NewExponentialBackOff(), 30))
	if err != nil {
		return err
	}
	if directory == "" {
		return nil
	}
	return migrateDB(edb.DB, directory)
}

type Person struct {
}

func (edb *ExampleDB) GetPersonByID(ctx context.Context, id int) (*Person, error) {
	panic("not implemented")
}

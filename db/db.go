package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cenk/backoff"
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
	"time"
)

func migrateDB(db *sqlx.DB, directory string) error {
	driver, _ := mysql.WithInstance(db.DB, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+directory,
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

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
	*sqlx.DB
}

func NewExampleDB(dsn string) (*ExampleDB, error) {
	return &ExampleDB{dsn, nil}, nil
}

func (edb *ExampleDB) ConnectAndMigrate(directory string) error {
	err := backoff.Retry(func() error {
		conn, err := sqlx.Open("mysql", edb.dsn)
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
	ID    int64
	Name  string
	Birth time.Time
	Death time.Time
}

func (edb *ExampleDB) GetPersonByID(ctx context.Context, id int) (*Person, error) {
	person := &Person{}
	err := edb.GetContext(ctx, person, "SELECT * FROM `person` WHERE `id` = ?", id)
	return person, err
}

func (edb *ExampleDB) AddPerson(ctx context.Context, p *Person) error {
	result := edb.MustExec("INSERT into `person` (`name`, `birth`, `death`) VALUES ( ?, ?, ?)", p.Name, p.Birth, p.Death)
	return checkResultAndSetID(&p.ID, result)
}

func checkResultAndSetID(id *int64, result sql.Result) error {
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("insert did not add a row")
	}
	*id, err = result.LastInsertId()
	return err
}

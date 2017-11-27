package example

import (
	"context"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_example.go -package=mocks github.com/bryanpaluch/example_go_app/example DB

type Person struct {
	ID    int64      `json:"id" db:"id"`
	Name  string     `json:"name" db:"name"`
	Birth *time.Time `json:"birth,omitempty" db:"birth"`
	Death *time.Time `json:"death,omitempty" db:"death"`
}

type DB interface {
	GetPersonByID(ctx context.Context, id int) (*Person, error)
	AddPerson(ctx context.Context, p *Person) error
}

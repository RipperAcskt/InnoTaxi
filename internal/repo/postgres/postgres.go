package postgres

import (
	"database/sql"
	"fmt"

	"github.com/RipperAcskt/innotaxi/internal/model"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrUserAlreadyExists error = fmt.Errorf("user already exists")
)

type Postgres struct {
	db      *sql.DB
	Migrate *migrate.Migrate
}

func New(url string) (*Postgres, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("open failed: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("with database failed")
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/repo/migrations", "postgres", driver)
	if err != nil {

		return nil, fmt.Errorf("new with database instance failed: %v", err)
	}

	return &Postgres{
		db,
		m,
	}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) CreateUser(user model.UserSingUp) error {
	rows, err := p.db.Query("SELECT * FROM users")
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}

	defer rows.Close()

	if rows.Next() {
		return ErrUserAlreadyExists
	}

	_, err = p.db.Exec("INSERT INTO users (name, phone_number, email, password, raiting) VALUES($1, $2, $3, $4, 4.0)", user.Name, user.PhoneNumber, user.Email, []byte(user.Password))
	if err != nil {
		return fmt.Errorf("exec failed: %v", err)
	}
	return nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/RipperAcskt/innotaxi/internal/service"
	"github.com/pkg/errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	db      *sql.DB
	Migrate *migrate.Migrate
}

func New(url string) (*Postgres, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("with instance failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(os.Getenv("MIGRATEPATH"), "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("new with database instance failed: %w", err)
	}

	return &Postgres{
		db,
		m,
	}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) CreateUser(ctx context.Context, user service.UserSingUp) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row, err := p.db.QueryContext(queryCtx, "SELECT * FROM users WHERE phone_number = $1 OR email = $2", user.PhoneNumber, user.Email)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if row.Next() {
		return errors.Wrapf(service.ErrUserAlreadyExists, "user: %v", user.Name)
	}

	_, err = p.db.Exec("INSERT INTO users (name, phone_number, email, password, raiting) VALUES($1, $2, $3, $4, 4.0)", user.Name, user.PhoneNumber, user.Email, []byte(user.Password))
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (p *Postgres) CheckUserByEmail(ctx context.Context, email string) (*service.UserSingIn, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row, err := p.db.QueryContext(queryCtx, "SELECT phone_number, password FROM users WHERE phone_number = $1", email)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	if row.Next() {
		var user service.UserSingIn
		err := row.Scan(&user.PhoneNumber, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		return &user, nil
	}

	return nil, service.ErrUserDoesNotExists
}

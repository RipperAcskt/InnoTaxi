package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	db      *sql.DB
	Migrate *migrate.Migrate
	cfg     *config.Config
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := sql.Open("pgx", cfg.GetDBUrl())
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("with instance failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.MIGRATE_PATH, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("new with database instance failed: %w", err)
	}

	return &Postgres{
		db,
		m,
		cfg,
	}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) CreateUser(ctx context.Context, user service.UserSingUp) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.db.QueryRowContext(queryCtx, "SELECT name FROM users WHERE phone_number = $1 OR email = $2", user.PhoneNumber, user.Email).Scan(&user.Name)
	if err == nil {
		return fmt.Errorf("user: %v: %w", user.Name, service.ErrUserAlreadyExists)

	}

	_, err = p.db.ExecContext(ctx, "INSERT INTO users (name, phone_number, email, password, raiting) VALUES($1, $2, $3, $4, 4.0)", user.Name, user.PhoneNumber, user.Email, []byte(user.Password))
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (p *Postgres) CheckUserByPhoneNumber(ctx context.Context, phone_number string) (*service.UserSingIn, uint64, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := p.db.QueryRowContext(queryCtx, "SELECT id, phone_number, password FROM users WHERE phone_number = $1", phone_number)

	var id uint64
	var user service.UserSingIn

	err := row.Scan(&id, &user.PhoneNumber, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, service.ErrUserDoesNotExists
		}

		return nil, 0, fmt.Errorf("scan failed: %w", err)
	}

	return &user, id, nil
}

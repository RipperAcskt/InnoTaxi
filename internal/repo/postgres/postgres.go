package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/RipperAcskt/innotaxi/internal/model"
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
	db, err := sql.Open("pgx", cfg.GetPostgresUrl())
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

	var name string
	err := p.db.QueryRowContext(queryCtx, "SELECT name FROM users WHERE (phone_number = $1 OR email = $2) AND status = $3", user.PhoneNumber, user.Email, model.StatusCreated).Scan(&name)
	if err == nil {
		return fmt.Errorf("user: %v: %w", user.Name, service.ErrUserAlreadyExists)

	}

	_, err = p.db.ExecContext(ctx, "INSERT INTO users (name, phone_number, email, password, raiting, status) VALUES($1, $2, $3, $4, 0.0, $5)", user.Name, user.PhoneNumber, user.Email, []byte(user.Password), model.StatusCreated)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (p *Postgres) CheckUserByPhoneNumber(ctx context.Context, phone_number string) (*service.UserSingIn, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := p.db.QueryRowContext(queryCtx, "SELECT id, phone_number, password FROM users WHERE phone_number = $1 AND status = $2", phone_number, model.StatusCreated)

	var user service.UserSingIn

	err := row.Scan(&user.ID, &user.PhoneNumber, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUserDoesNotExists
		}

		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return &user, nil
}

func (p *Postgres) GetUserById(ctx context.Context, id string) (*model.User, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &model.User{}
	err := p.db.QueryRowContext(queryCtx, "SELECT id, name, phone_number, email, raiting FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Email, &user.Raiting)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUserDoesNotExists
		}
		return nil, fmt.Errorf("query row context failed: %w", err)
	}

	return user, err
}

func (p *Postgres) UpdateUserById(ctx context.Context, id string, user *model.User) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := p.db.ExecContext(queryCtx, "UPDATE users SET name = COALESCE($1, name), phone_number = COALESCE($2, phone_number), email = COALESCE($3, email) WHERE id = $4 AND status = $5", user.Name, user.PhoneNumber, user.Email, id, model.StatusCreated)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return service.ErrUserDoesNotExists
	}
	return nil
}

func (p *Postgres) DeleteUserById(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := p.db.ExecContext(queryCtx, "UPDATE users SET status = $1 WHERE id = $2 AND status = $3", model.StatusDeleted, id, model.StatusCreated)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return service.ErrUserDoesNotExists
	}
	return nil
}

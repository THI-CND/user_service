package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/BieggerM/userservice/pkg/models"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database interface {
	Connect(dbHost, dbPort, dbUser, dbPassword, dbName string) error
	SaveUser(user models.User) error
	DeleteUser(username string)
	UpdateUser(user models.User) (models.User, error)
	GetUser(username string) (models.User, error)
	ListUsers() []models.User
	RunMigrations(migrationPath string) error
	Close() error
}

// Postgres is the PostgreSQL database connection
type Postgres struct {
	DB *sql.DB
}

// Connect connects to the PostgreSQL database
func (p *Postgres) Connect(dbHost, dbPort, dbUser, dbPassword, dbName string) error {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	p.DB, err = sql.Open("postgres", connStr)
	err = p.DB.Ping()
	return err
}

// SaveUser saves a user to the PostgreSQL database
func (p *Postgres) SaveUser(user models.User) error {
	var exists bool
	err := p.DB.QueryRow("select exists(select 1 from users where username=$1)", user.Username).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}
	_, err = p.DB.Exec("insert into users (username, firstname, lastname, password) values ($1, $2, $3, $4)", user.Username, user.FirstName, user.LastName, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user from the PostgreSQL database
func (p *Postgres) DeleteUser(username string) {
	_, err := p.DB.Exec("delete from users where username = $1", username)
	if err != nil {
		fmt.Println(err)
	}
}

// UpdateUser updates a user in the PostgreSQL database
func (p *Postgres) UpdateUser(user models.User) (models.User, error) {
	res, err := p.DB.Exec("update users set firstname = $1, lastname = $2 where username = $3", user.FirstName, user.LastName, user.Username)
	if err != nil {
		return user, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return user, err
	}
	if rowsAffected == 0 {
		return user, errors.New("user does not exist")
	}
	return user, nil
}

// GetUser gets a user from the PostgreSQL database
func (p *Postgres) GetUser(username string) (models.User, error) {
	user := models.User{}
	err := p.DB.QueryRow("select username, firstname, lastname, password from users where username = $1", username).Scan(&user.Username, &user.FirstName, &user.LastName, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No user found with the given username")
			return user, err
		} else {
			fmt.Println(err)
		}
	}
	return user, nil
}

// ListUsers lists all users from the PostgreSQL database
func (p *Postgres) ListUsers() []models.User {
	rows, err := p.DB.Query("select username, firstname, lastname from users")
	if err != nil {
		fmt.Println(err)
	}
	var users []models.User
	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.Username, &user.FirstName, &user.LastName)
		users = append(users, user)
	}
	return users
}

// RunMigrations runs the database migrations
func (p *Postgres) RunMigrations(migrationPath string) error {
	driver, err := postgres.WithInstance(p.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run up migrations: %w", err)
	}
	return nil
}

func (p *Postgres) Close() error {
	return p.DB.Close()
}

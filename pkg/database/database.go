package database

import (
    "database/sql"
    "fmt"
    "github.com/BieggerM/userservice/pkg/models"
    _ "github.com/lib/pq"
)

type Database interface {
    SaveUser(user models.User)
    DeleteUser(username string)
    UpdateUser(user models.User)
    GetUser(username string) models.User
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
    return err
}

// SaveUser saves a user to the PostgreSQL database
func (p *Postgres) SaveUser(user models.User) error {
    _, err := p.DB.Exec("insert into users (username, firstname, lastname) values ($1, $2, $3)", user.Username, user.FirstName, user.LastName)
    if err != nil {
        fmt.Println(err)
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
func (p *Postgres) UpdateUser(user models.User) models.User {
    _, err := p.DB.Exec("update users set firstname = $1, lastname = $2 where username = $3", user.FirstName, user.LastName, user.Username)
    if err != nil {
        fmt.Println(err)
    }
    return user
}

// GetUser gets a user from the PostgreSQL database
func (p *Postgres) GetUser(username string) models.User {

    user := models.User{}
    err := p.DB.QueryRow("select username, firstname, lastname from users where username = $1", username).Scan(&user.Username, &user.FirstName, &user.LastName)
    if err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("No user found with the given username")
        } else {
            fmt.Println(err)
        }
    }
    return user
}

// ListUsers lists all users from the PostgreSQL database
func (p *Postgres) ListUsers() []models.User {
    rows, err := p.DB.Query("select * from users")
    if err != nil {
        fmt.Println(err)
    }
    users := []models.User{}
    for rows.Next() {
        user := models.User{}
        rows.Scan(&user.Username, &user.FirstName, &user.LastName)
        users = append(users, user)
    }
    return users
}

func (p *Postgres) Close() error {
    return p.DB.Close()
}
package services

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	config "github.com/sriharivishnu/shopify-challenge/config"
)

var userTable = `
create table if not exists user (
    id CHAR(36) PRIMARY KEY default uuid(),
    username VARCHAR(256) not null UNIQUE,
    password VARCHAR(256) not null,
    created_at timestamp default now()
);
`

var repositoryTable = `
create table if not exists repository (
    id CHAR(36) PRIMARY KEY default uuid(),
	owner_id CHAR(36) NOT NULL,
	name TEXT,
	description VARCHAR(2056),
    created_at timestamp default now(),
    INDEX(owner_id),
    FULLTEXT(name),
	UNIQUE(owner_id, name),
    FOREIGN KEY (owner_id) REFERENCES user(id)
    ON DELETE CASCADE
);
`
var imageTagTable = `
create table if not exists image_tag (
    id CHAR(36) PRIMARY KEY default uuid(),
	repository_id CHAR(36) NOT NULL,
    tag VARCHAR(256) NOT NULL,
    description VARCHAR(2056),
	file_key TEXT NOT NULL,
    created_at timestamp default now(),
    INDEX (repository_id),
	UNIQUE(repository_id, tag),
    FOREIGN KEY (repository_id) REFERENCES repository(id)
		ON DELETE CASCADE
);
`

var DbConn *sqlx.DB

func createSchema() {
	log.Println("Creating Database...")
	DbConn.MustExec(fmt.Sprintf("create database if not exists %s;", config.Config.DATABASE_NAME))
	log.Println("Created Database")
	DbConn.MustExec(fmt.Sprintf("use %s;", config.Config.DATABASE_NAME))
	log.Println("Creating Schema...")
	DbConn.MustExec(userTable)
	DbConn.MustExec(repositoryTable)
	DbConn.MustExec(imageTagTable)
	log.Println("Done creating schema")
}

func Init() {
	host := config.Config.DATABASE_HOST
	port := config.Config.DATABASE_PORT
	user := config.Config.DATABASE_USER
	pass := config.Config.DATABASE_PASSWORD
	DbConn = sqlx.MustConnect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/?parseTime=true", user, pass, host, port))
	log.Println("Connected to DB host")
	log.Println("Verifying schema...")
	createSchema()
	log.Println("Database is connected and ready!")
}

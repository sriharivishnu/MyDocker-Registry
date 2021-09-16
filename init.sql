create database if not exists prod;
use prod;

create table if not exists user (
    id CHAR(36) PRIMARY KEY default uuid(),
    username VARCHAR(256) not null UNIQUE,
    password VARCHAR(256) not null,
    created_at timestamp default now()
);

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
CREATE TABLE IF NOT EXISTS users (
    id int NOT NULL AUTO_INCREMENT,
    full_name varchar(50) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    password varchar(72) NOT NULL,
    is_admin tinyint NOT NULL DEFAULT 0,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW(),
    deleted_at datetime NULL,

    PRIMARY KEY(id) 
)
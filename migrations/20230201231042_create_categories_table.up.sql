CREATE TABLE categories (
    id int NOT NULL AUTO_INCREMENT,

    uuid char(36) NOT NULL,
    parent_id int NULL,
    title varchar(50) NOT NULL,
    description  varchar(255) NULl DEFAULT NULL,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW(),
    deleted_at datetime NULL,

    PRIMARY KEY (id)
)
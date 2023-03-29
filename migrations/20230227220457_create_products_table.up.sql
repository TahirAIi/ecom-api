CREATE TABLE IF NOT EXISTS products (
    id int NOT NULL AUTO_INCREMENT,
    category_id int NULL,
    title varchar(50) NOT NULL,
    description  varchar(255) NULl DEFAULT NULL,
    price int NOT NULL,
    main_picture varchar(50) NULl,
    created_at datetime DEFAULT NOW(),
    updated_at datetime DEFAULT NOW(),
    deleted_at datetime NULL,

    PRIMARY KEY (id),
    FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
)
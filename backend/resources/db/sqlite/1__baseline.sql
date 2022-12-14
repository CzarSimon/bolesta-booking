-- +migrate Up
CREATE TABLE IF NOT EXISTS `cabin` (
    `id` VARCHAR(50) NOT NULL,  
    `name` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (id),
    UNIQUE(name)
);
CREATE TABLE IF NOT EXISTS `user_account` (
    `id` VARCHAR(50) NOT NULL,  
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(256) NOT NULL,
    `salt` VARCHAR(64) NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    PRIMARY KEY (id),
    UNIQUE(`name`),
    UNIQUE(`email`),
    UNIQUE(`password`),
    UNIQUE(`salt`)
);
CREATE TABLE IF NOT EXISTS `booking` (
    `id` VARCHAR(50) NOT NULL,  
    `start_date` DATETIME NOT NULL,
    `end_date` DATETIME NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    `cabin_id` VARCHAR(50) NOT NULL,
    `user_id` VARCHAR(50) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (cabin_id) REFERENCES cabin (id),
    FOREIGN KEY (user_id) REFERENCES user_account (id)
);
INSERT OR IGNORE INTO cabin(id, name, created_at, updated_at) VALUES
    ('a4b4f496-767e-423e-9816-83b71e1cfa89', 'Bölestastugan', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('63e71fef-0037-451f-b731-27249c0164d9', 'Gulhuset', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('2aa15162-2443-48f1-9b8f-6314f90faf9a', 'Bergebo', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
-- +migrate Down
DROP TABLE IF EXISTS `booking`;
DROP TABLE IF EXISTS `user_account`;
DROP TABLE IF EXISTS `cabin`;
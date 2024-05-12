CREATE TABLE mailbox
(
    username VARCHAR(254) PRIMARY KEY UNIQUE NOT NULL,
    password VARCHAR(256)                    NOT NULL,
    active   INT                             NOT NULL
);

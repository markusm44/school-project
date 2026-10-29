CREATE TABLE IF NOT EXISTS `users` (
    `userID`      INTEGER PRIMARY KEY AUTOINCREMENT,
    `nickname`    TEXT        NOT NULL,
    `email`       TEXT UNIQUE NOT NULL,
    `firstName`   TEXT        NOT NULL,
    `lastName`    TEXT        NOT NULL,
    `password`    TEXT        NOT NULL,
    `image`       TEXT        NOT NULL DEFAULT 'default.png',
    `isProfilePublic` BOOLEAN DEFAULT 1,
    `dateOfBirth` TEXT,
    `aboutMe`     TEXT
);
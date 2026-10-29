CREATE TABLE IF NOT EXISTS `groups` (
    `groupID`  INTEGER PRIMARY KEY AUTOINCREMENT,
    `creatorID` INTEGER NOT NULL,
    `name`       TEXT NOT NULL,
    `description` TEXT,
    `members` BLOB,
    `created` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(creatorID) REFERENCES users(userID) ON DELETE CASCADE
);
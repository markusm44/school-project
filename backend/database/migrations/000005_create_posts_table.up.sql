CREATE TABLE IF NOT EXISTS `posts` (
    `postID`      INTEGER PRIMARY KEY AUTOINCREMENT,
    `userID`      INTEGER NOT NULL,
    `content`     TEXT    NOT NULL,
    `created`     TEXT    NOT NULL,
    `image`       TEXT    NOT NULL DEFAULT 'default.png',
    `fullName`    TEXT    NOT NULL, 
    `status`      TEXT    NOT NULL,
    `almostPrivateUsers`  VARBINARY NOT NULL, 
    FOREIGN KEY (userID) REFERENCES users (userID) ON DELETE CASCADE
);

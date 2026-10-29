CREATE TABLE IF NOT EXISTS `groupPosts` (
    `postID`      INTEGER PRIMARY KEY AUTOINCREMENT,
    `userID`      INTEGER NOT NULL,
    `groupID`     INTEGER NOT NULL,
    `content`     TEXT    NOT NULL,
    `created`     TEXT    NOT NULL,
    `image`       TEXT    NOT NULL DEFAULT 'default.png',
    `fullName`    TEXT    NOT NULL, 
    FOREIGN KEY (userID) REFERENCES users (userID) ON DELETE CASCADE
);

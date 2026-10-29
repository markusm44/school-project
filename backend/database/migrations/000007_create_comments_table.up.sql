CREATE TABLE IF NOT EXISTS `comments`(
    `commentID`  INTEGER PRIMARY KEY AUTOINCREMENT,
    `postID`     INTEGER NOT NULL,
    `content`    TEXT NOT NULL,
    `userID`     INTEGER NOT NULL,
    `fullName`   TEXT    NOT NULL,
    `created`    TEXT    NOT NULL,
    `image`      TEXT    NOT NULL DEFAULT 'default_image.png',
 FOREIGN KEY (userID) REFERENCES users (userID) ON DELETE CASCADE,
 FOREIGN KEY (postID) REFERENCES posts (postID) ON DELETE CASCADE
);
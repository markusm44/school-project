CREATE TABLE IF NOT EXISTS `notifications` (
    `notificationID`     INTEGER PRIMARY KEY AUTOINCREMENT,
    `notificationSubtype`   TEXT,
    `senderID`           INTEGER,
    `recieverID`         INTEGER,
    `content`            TEXT NOT NULL,
    `created`            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (senderID) REFERENCES users (userID) ON DELETE CASCADE
  FOREIGN KEY (recieverID) REFERENCES users (userID) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS `events` (
    `eventID`     INTEGER PRIMARY KEY AUTOINCREMENT,
    `userID`      INTEGER NOT NULL,
    `groupID`     INTEGER NOT NULL,
    `content`     TEXT    NOT NULL,
    `created`     TEXT    NOT NULL,
    `title`       TEXT,
    `date`        TEXT    NOT NULL,
  FOREIGN KEY (userID) REFERENCES users (userID) ON DELETE CASCADE
);
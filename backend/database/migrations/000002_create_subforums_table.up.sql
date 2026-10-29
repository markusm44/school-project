CREATE TABLE IF NOT EXISTS `subforums` (
    `subforumID`  INTEGER PRIMARY KEY AUTOINCREMENT,
    `title`       TEXT NOT NULL,
    `description` TEXT,
    `created` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS targets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    mission_id INTEGER,
    name TEXT,
    country TEXT,
    notes TEXT,
    complete BOOLEAN,
    FOREIGN KEY (mission_id) REFERENCES missions(id)
);

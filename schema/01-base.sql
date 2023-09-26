-- Base schema for cave-logger

-- trip
--------------------------------------------------------------------------------
DROP TABLE IF EXISTS trip;

CREATE TABLE trip (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date INTEGER NOT NULL,
    caveid INTEGER NOT NULL,
    notes BLOB,
    FOREIGN KEY (caveid) REFERENCES cave (id)
);

DROP TRIGGER IF EXISTS trips_stats_insert;

CREATE TRIGGER IF NOT EXISTS trip_stats_insert
    AFTER INSERT ON trip
    FOR EACH ROW
BEGIN
    INSERT INTO stats (kind, value, count)
    VALUES (
        'cave',
        NEW.caveid,
        (
            SELECT IFNULL(COUNT(*), 0) + 1
            FROM stats
            WHERE kind = 'cave' AND value = NEW.caveid
        )
    );
END;

DROP TRIGGER IF EXISTS trips_stats_delete;

CREATE TRIGGER IF NOT EXISTS trip_stats_delete
    AFTER DELETE ON trip
    FOR EACH ROW
BEGIN
    INSERT INTO stats (kind, value, count)
    VALUES (
        'cave',
        OLD.caveid,
        (
            SELECT IFNULL(COUNT(*)-1, 0)
            FROM stats
            WHERE kind = 'cave' AND value = OLD.caveid
        )
    );
END;

-- trip_group
--------------------------------------------------------------------------------
DROP TABLE IF EXISTS trip_group;

CREATE TABLE trip_group (
    tripid INTEGER NOT NULL,
    caverid INTEGER NOT NULL,
    FOREIGN KEY (tripid) REFERENCES trip (id)
    FOREIGN KEY (caverid) REFERENCES caver (id)
);

DROP TRIGGER IF EXISTS trip_group_stats_insert;

CREATE TRIGGER IF NOT EXISTS trip_group_stats_insert
    AFTER INSERT ON trip_group
    FOR EACH ROW
BEGIN
    INSERT INTO stats (kind, value, count)
    VALUES (
        'caver',
        NEW.caverid,
        (
            SELECT IFNULL(COUNT(*), 0) + 1
            FROM stats
            WHERE kind = 'caver' AND value = NEW.caverid
        )
    );
END;

DROP TRIGGER IF EXISTS trip_group_stats_delete;

CREATE TRIGGER IF NOT EXISTS trip_group_stats_delete
    AFTER DELETE ON trip_group
    FOR EACH ROW
BEGIN
    INSERT INTO stats (kind, value, count)
    VALUES (
        'caver',
        OLD.caverid,
        (
            SELECT IFNULL(COUNT(*), 0) - 1
            FROM stats
            WHERE kind = 'caver' AND value = OLD.caverid
        )
    );
END;

--------------------------------------------------------------------------------
DROP TABLE IF EXISTS cave;

CREATE TABLE cave (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    region TEXT,
    country TEXT,
    srt INTEGER NOT NULL,
    notes BLOB
);

--------------------------------------------------------------------------------
DROP TABLE IF EXISTS caver;

CREATE TABLE caver (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    club TEXT,
    notes BLOB
);

--------------------------------------------------------------------------------
DROP TABLE IF EXISTS stats;

CREATE TABLE stats (
    kind TEXT NOT NULL,
    value INTEGER NOT NULL,
    count INTEGER NOT NULL
);

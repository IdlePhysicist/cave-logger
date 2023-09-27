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


DROP TRIGGER IF EXISTS trips_stats_delete;


-- trip_group
--------------------------------------------------------------------------------
DROP TABLE IF EXISTS trip_group;

CREATE TABLE trip_group (
    tripid INTEGER NOT NULL,
    caverid INTEGER NOT NULL,
    FOREIGN KEY (tripid) REFERENCES trip (id)
    FOREIGN KEY (caverid) REFERENCES caver (id)
    UNIQUE (tripid, caverid)
);

DROP TRIGGER IF EXISTS trip_group_stats_insert;


DROP TRIGGER IF EXISTS trip_group_stats_delete;


--------------------------------------------------------------------------------
DROP TABLE IF EXISTS cave;

CREATE TABLE cave (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    region TEXT,
    country TEXT,
    srt INTEGER NOT NULL,
    notes BLOB,
    UNIQUE (name, region, country)
);

--------------------------------------------------------------------------------
DROP TABLE IF EXISTS caver;

CREATE TABLE caver (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    club TEXT,
    notes BLOB,
    UNIQUE (name, club)
);

--------------------------------------------------------------------------------
DROP TABLE IF EXISTS stats;

CREATE TABLE stats (
    kind TEXT NOT NULL,
    value INTEGER NOT NULL,
    count INTEGER NOT NULL,
    UNIQUE (kind, value)
);

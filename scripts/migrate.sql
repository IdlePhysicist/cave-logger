-- This script is to migrate from the old schema to the new one

INSERT INTO caver (id, name, club, notes)
    SELECT * FROM people;

INSERT INTO cave (id, name, region, country, srt, notes)
    SELECT * FROM locations;

INSERT INTO trip (id, date, caveid, notes)
    SELECT t.id, date(t.date, 'unixepoch', 'localtime'), t.caveid, t.notes
    FROM trips t;

INSERT INTO trip_group (tripid, caverid)
    SELECT * FROM trip_groups;


-- Finally drop the old tables
-- DROP TABLE IF EXISTS trips;
-- DROP TABLE IF EXISTS trip_groups;
-- DROP TABLE IF EXISTS people;
-- DROP TABLE IF EXISTS locations;
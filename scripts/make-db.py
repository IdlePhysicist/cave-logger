#!/usr/bin/env python3
from datetime import datetime
import os
import sys
import sqlite3
import uuid
import json

NOW = datetime.now().strftime("%Y-%m-%dT%H_%M")
HOME = os.environ["HOME"]
NEWPATH = f"{HOME}/.config/cave-logger"

def runTableMigration(dbFile):
    """
    This function migrates old databases to the new spec.
    """
    print(f"Checking existing database for new spec")

    conn = sqlite3.connect(dbFile)
    c = conn.cursor()

    def check(cols):
        for col_info in cols:
            if "notes" in col_info:
                return True
        return False
    #END

    for t in ("locations","people"):
        c.execute(f"PRAGMA table_info({t});")
        if not check(c.fetchall()):
            stmt = f"ALTER TABLE {t} ADD notes TEXT;"
            c.execute(stmt)
            print(f"Added notes column to {t}")

    conn.commit()
    conn.close()
#END

try:
  os.makedirs(NEWPATH, 0o755)
except FileExistsError:
  print("Config directory exists \nChecking for existing configs")

  # Now we check the config/cave-logger dir for files that we would expect
  os.chdir(NEWPATH)
  files = os.listdir()
  if files:
    db = False
    cfg = False
    runMigration = False
    for f in files:
      if '.db' in f:
        db = f
      elif 'config.json' == f:
        cfg = f

    with open(cfg, 'r') as f:
      try:
        fn = json.load(f)['database']['filename']
        if db in fn:
          print("Found pre-existing configs aborting...")
          runMigration = True
      except KeyError:
        pass

    if runMigration:
        runTableMigration(f"{HOME}/{fn}")
    sys.exit()

sqliteFile = str(uuid.uuid4()) + '.db'
tables = [
  {
    "name" : "trips",
    "cols" : ["id","date","caveid","notes"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","INTEGER","INTEGER","TEXT"]
  },
  {
    "name" : "locations",
    "cols" : ["id","name","region","country","srt","notes"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","TEXT","TEXT","TEXT","INTEGER","TEXT"]
  },
  {
    "name" : "people",
    "cols" : ["id","name","club","notes"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","TEXT","TEXT","TEXT"]
  },
  {
    "name" : "trip_groups",
    "cols" : ["tripid","caverid"],
    "types": ["INTEGER","INTEGER"]
  }
]

conn = sqlite3.connect("/".join([NEWPATH, sqliteFile]))
print("Connected to new sqlite database")
print("Database file name: ", sqliteFile)

c = conn.cursor()

for table in tables:
  fields = ','.join(
    [ 
      ' '.join(
        [ table['cols'][i], table['types'][i] ]
      ) 
      for i in range(len(table['cols']))
    ]
  )

  stmt = 'CREATE TABLE {tn} ({flds})'.format(
      tn=table['name'], flds=fields
    )

  c.execute(stmt)

conn.commit()
print("Created {} database tables".format(len(tables)))
conn.close()


CONFIG_FN = '{}/config.json'.format(NEWPATH)
with open(CONFIG_FN, 'w') as c:
  config = {
    'database': {
      'filename': '/'.join([".config/cave-logger", sqliteFile]),
      'created' : NOW,
      'spec': "1.1.0"
    },
    'colors': {
      'primitiveBackground': '',
      'contrastBackground': '',
      'moreContrastBackground': '',
      'border': '',
      'title': '',
      'graphics': '',
      'primaryText': '',
      'secondaryText': '',
      'tertiaryText': '',
      'inverseText': '',
      'contrastSecondaryText': ''
    }
  }
  json.dump(config, c, indent=2)
  print("Wrote database name to config file")


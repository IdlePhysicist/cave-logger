#!/usr/bin/env python
from datetime import datetime
import os
import sys
import sqlite3
import uuid
import yaml

NOW = datetime.now().strftime("%Y-%m-%dT%H_%M")
HOME = os.environ["HOME"]
NEWPATH = "{}/.config/cave-logger".format(HOME)

if sys.version_info < (3,):
  print("""
Why are you not using Python 3 by now?
Even I am...
Alas this script will still try to run.\n"""
  )


try:
  os.makedirs(NEWPATH, 0o755)
except FileExistsError:
  print("Directory exists moving on")


sqliteFile = str(uuid.uuid4()) + '.db'
tables = [
  {
    "name" : "trips",
    "cols" : ["id","date","caveid","notes"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","INTEGER","INTEGER","INTEGER","TEXT"]
  },
  {
    "name" : "locations",
    "cols" : ["id","name","region","country","srt"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","TEXT","TEXT","TEXT","INTEGER"]
  },
  {
    "name" : "people",
    "cols" : ["id","name","club"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","TEXT","TEXT"]
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

CONFIG_FN = '{}/config.yml'.format(NEWPATH)
with open(CONFIG_FN, 'w') as c:
  config = {
    'database': {
      'filename': '/'.join([NEWPATH, sqliteFile]),
      'created' : NOW
    }
  }
  yaml.dump(config, c)
  print("Wrote database name to config file")

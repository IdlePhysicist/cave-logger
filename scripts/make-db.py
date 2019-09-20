#!/usr/bin/env python
from datetime import datetime
import sys
import sqlite3
import uuid
import yaml

NOW = datetime.now().strftime("%Y-%m-%dT%H_%M")

if sys.version_info < (3,):
  print("""
Why are you not using Python 3 by now?
Even I am...
Alas this script will still try to run.\n"""
  )

sqliteFile = str(uuid.uuid4()) + '.db'
tables = [
  {
    "name" : "entries",
    "cols" : ["id","date","caveid","caverids","notes"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","INTEGER","INTEGER","INTEGER","INTEGER","TEXT"]
  },
  {
    "name" : "caves",
    "cols" : ["id","name","region","country","srt","visits"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","INTEGER","TEXT","TEXT","INTEGER","INTEGER"]
  },
  {
    "name" : "cavers",
    "cols" : ["id","first", "last","club","freq"],
    "types": ["INTEGER PRIMARY KEY AUTOINCREMENT","TEXT","TEXT","TEXT","INTEGER"]
  }
]

conn = sqlite3.connect("/".join(["./config", sqliteFile]))
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

CONFIG_FN = './config/config.yml'
with open(CONFIG_FN, 'w') as c:
  config = {
    'database': {
      'filename': './config/{}'.format(sqliteFile),
      'created' : NOW
    }
  }
  yaml.dump(config, c)
  print("Wrote database name to config file")

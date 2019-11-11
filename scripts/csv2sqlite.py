#!/usr/bin/env python3

import csv
import json
import subprocess
import sqlite3
import yaml
import os
import sys

from datetime import datetime

HOME = os.environ['HOME']
#NAMES = f"{HOME}/Desktop/names.json"
CSV = sys.argv[1]

INITSCRIPT = "./scripts/make-db.py"
CFG = "./config/config.yml"

class App(object):
  def __init__(self):
    out = subprocess.run(INITSCRIPT, capture_output=True)
    
    with open(CFG, 'r') as c:
      self.config = yaml.safe_load(c)

    self.conn = sqlite3.connect(
      self.config['database']['filename'],
      check_same_thread=False
    )
    self.conn.row_factory = sqlite3.Row

  def read(self, inputFile):
    inputData = []
    with open(inputFile, 'r') as c:
      reader = csv.reader(c)
      lines = 0
      for row in reader:
        lines += 1
        if lines == 1:
          continue
        names = []
        for field in row[2:11]:
          if field == '':
            break
          names.append(field)
        inputData.append({
          'date': row[0],
          'cave': row[1],
          'names': names,
          'notes': row[11]
        })
    #print(inputData)
    return inputData

  def main(self):
    rows = self.read(CSV)

    for row in rows:
      code = self.process(row)
      print("SUCCESS")

    print("DONE")

  def process(self, row):
    entry = {}
    # Date
    stamp = int(datetime.strptime(row['date']+'T12:00:00Z', '%d/%m/%yT%H:%M:%SZ').timestamp())
    entry.update({'date': stamp})

    # Cave
    knownCaves = self.getAllCaves()
    if row['cave'] in knownCaves.keys():
      caveID = knownCaves[row['cave']]
    else:
      # Create a new cave
      caveID = self.createCave({'name': row['cave']})
    entry.update({'caveid': caveID})

    # Cavers
    knownCavers = self.getAllCavers()
    caverIDs = []
    for name in row['names']:
      if name in knownCavers.keys():
        caverIDs.append(str(knownCavers[name]))
      else:
        # Create a new one
        newCaverID = self.createCaver({'name': name})
        caverIDs.append(str(newCaverID))

    # Notes
    entry.update({'notes': row['notes']})

    tripID = self.createEntry(entry)

    result = self.createTripGroup(tripID, caverIDs)

  #
  # -- Creator Methods
  #
  def createCave(self, cave):
    query = "INSERT INTO locations (name) VALUES (?)"
    return self.insert(query, [ cave['name'] ])

  def createCaver(self, caver):
    query = "INSERT INTO people (name) VALUES (?)"
    return self.insert(query, [ caver['name'] ])

  def createEntry(self, trip):
    query = """
      INSERT INTO trips (date, caveid, notes)
      VALUES (?,?,?)"""
    return self.insert(
        query,
        [ trip['date'], trip['caveid'], trip['notes'] ]
      )

  def createTripGroup(self, trip, group):
    placeholders = []
    params = []
    for person in group:
      # Add to the placeholder list
      placeholders.append("(?,?)")

      # Add to the params list
      params.extend([trip, person])

    placeholder = ",".join(placeholders)
    query = f"""
      INSERT INTO trip_groups (tripid, caverid)
      VALUES {placeholder}"""

    return self.insert(query, params)

  #
  # -- Direct DB Methods
  #
  def insert(self, query, params):
    try:
      c = self.conn.cursor()
      c.execute(query, params)
      self.conn.commit()
      return c.lastrowid
    except sqlite3.Error as e:
      print(f"sqlite.insert: An error occured {e}")
      return []

  def getAllCavers(self):
    try:
      c = self.conn.cursor()
      query = "SELECT id, name FROM people"
      c.execute(query)
      rows = c.fetchall()
      if not rows:
        print(f"sqlite.get: Query returned empty")
        return {}
      return { r['name']: r['id'] for r in rows }

    except sqlite3.Error as e:
      print(f"sqlite.get: An error occured {e}")
      return {}

  def getAllCaves(self):
    try:
      c = self.conn.cursor()
      query = "SELECT id, name FROM locations"
      c.execute(query)
      rows = c.fetchall()
      if not rows:
        print("sqlite.get: Query returned empty")
        return {}
      return { r['name']: r['id'] for r in rows }
    except sqlite3.Error as e:
      print(f"sqlite.get: An error occured {e}")
      return {}


if __name__ == '__main__':
  App().main()


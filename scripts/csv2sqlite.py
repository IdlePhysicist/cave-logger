#!/usr/bin/env python3

import csv
import json
import subprocess
import sqlite3
import yaml
import os

from datetime import datetime

HOME = os.environ['HOME']
NAMES = f"{HOME}/Desktop/names.json"
CSV = f"{HOME}/Desktop/CaveLog.csv"

INITSCRIPT = "./scripts/make-db.py"
CFG = "./config/config.yml"

class App(object):
  def __init__(self):
    out = subprocess.run(INITSCRIPT, capture_output=True)
    print(out)

    # Open JSON & CSV
    with open(NAMES, 'r') as n:
      self.nameDict = json.load(n)
    
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
      if not code:
        print(f"FAILURE {row}")
        continue
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
    for caver in row['names']:
      if caver in self.nameDict.keys():
        fullCaver = '+'.join([self.nameDict[caver][0], self.nameDict[caver][1]])
      else:
        fullCaver = caver+'+'
      
      if fullCaver in knownCavers.keys():
        caverIDs.append(str(knownCavers[fullCaver]))
      else:
        # Create a new one
        if caver in self.nameDict.keys(): 
          first = self.nameDict[caver][0]
          last  = self.nameDict[caver][1]
        else:
          first = caver
          last = ''
        
        newCaverID = self.createCaver({
          'first': first,
          'last' : last
        })
        caverIDs.append(str(newCaverID))
    
    entry.update({'caverids': '|'.join(caverIDs)})
    
    # Notes
    entry.update({'notes': row['notes']})

    rowid = self.createEntry(entry)
    if not rowid:
      return False
    return True

  #
  # -- Creator Methods
  #
  def createCave(self, cave):
    query = "INSERT INTO caves (name) VALUES (?)"
    return self.insert(query, [ cave['name'] ])

  def createCaver(self, caver):
    query = "INSERT INTO cavers (first, last) VALUES (?,?)"
    return self.insert(query, [ caver['first'], caver['last'] ])

  def createEntry(self, entry):
    query = """
      INSERT INTO entries (date, caveid, caverids, notes)
      VALUES (?,?,?,?)"""
    return self.insert(
        query,
        [ entry['date'], entry['caveid'], entry['caverids'], entry['notes'] ]
      )

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
      query = "SELECT id, first, last FROM cavers"
      c.execute(query)
      rows = c.fetchall()
      if not rows:
        print(f"sqlite.get: Query returned empty")
        return {}
      return { '+'.join([r['first'],r['last']]): r['id'] for r in rows }

    except sqlite3.Error as e:
      print(f"sqlite.get: An error occured {e}")
      return {}

  def getAllCaves(self):
    try:
      c = self.conn.cursor()
      query = "SELECT id, name FROM caves"
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


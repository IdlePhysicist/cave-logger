#!/usr/bin/env python3
import argparse
import csv
import sqlite3
import yaml
from collections import Counter

CONFIG_FN = "./config/config.yml"

class Main(object):
  def __init__(self):
    with open(CONFIG_FN, 'r') as y:
      self.config = yaml.safe_load(y)

    self.conn = sqlite3.connect(self.config['database']['filename'])

  def main(self):
    inputData = self.parse()
    x = Counter([ r['cave'] for r in inputData ])
    caves = x.keys()
    caveFreq = x.values()
    uniqueCaves = [ {caves[i]:caveFreq[i]} for i in range(len(caves)) ]
    
    y = 

  def parse(self):
    parser = argparse.ArgumentParser()
    parser.add_argument('-i', type=str)
    args = parser.parse_args()
    return self.processCsv(args.i)

  def processCsv(self, inputFile):
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
    print(inputData)
    return inputData

  def insert(self, data):
    

if __name__ == '__main__':
  Main().main()

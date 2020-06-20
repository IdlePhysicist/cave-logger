#!/usr/bin/env python3

import json
import os
import shutil

HOME = os.environ["HOME"]
CFGFILE = "./config/config.json"

NEWPATH = ".config/cave-logger"

try:
  os.makedirs(NEWPATH, 0o755)
except FileExistsError:
  print("Directory exists moving on")

with open(CFGFILE, 'r') as c:
  cfg = json.load(c)

shutil.copy(
  '/'.join([HOME,cfg['database']['filename']]),
  f"{HOME}/{NEWPATH}/"
)
shutil.copy(CFGFILE, f"{HOME}/{NEWPATH}/")

cfg['database']['filename'] = f"{HOME}/{NEWPATH}/{cfg['database']['filename'].split('/')[-1]}"

with open(f"{HOME}/{NEWPATH}/config.json", 'w') as c:
  json.dump(cfg, c)

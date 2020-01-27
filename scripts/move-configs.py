#!/usr/bin/env python3

import yaml
import os
import shutil

HOME = os.environ["HOME"]
CFGFILE = "./config/config.yml"

NEWPATH = ".config/cave-logger"

try:
  os.makedirs(NEWPATH, 0o755)
except FileExistsError:
  print("Directory exists moving on")

with open(CFGFILE, 'r') as c:
  cfg = yaml.safe_load(c)

shutil.copy(
  '/'.join([HOME,cfg['database']['filename']]),
  f"{HOME}/{NEWPATH}/"
)
shutil.copy(CFGFILE, f"{HOME}/{NEWPATH}/")

cfg['database']['filename'] = f"{HOME}/{NEWPATH}/{cfg['database']['filename'].split('/')[-1]}"

with open(f"{HOME}/{NEWPATH}/config.yml", 'w') as c:
  yaml.dump(cfg, c)

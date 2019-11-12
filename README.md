<p align="center"><img alt="Cave Logger" src="assets/logo.png"></p>

## Summary
Cave Logger is a basic SQLite database interface written in Go, and it allows cavers to track the caves that they have been to, who with, and when. 

I indend to make the code more generic to allow other outdoorsy people to use this app with less fuss.

## What It Looks Like
<p align="center"><img src="assets/screenshot.png"></p>

## Getting started

### Install
You can install by the following set of instructions:

1. Clone or download the repo, and naviagte to the repo directory
2. Open the Makefile and modify the `GOOS` variable to your platform. Default is `darwin`
3. Following this run `make`, this will produce the binary `./build/cave-logger`
4. Assuming the binary has built correctly then you have two courses of action:
    - A. If you have no data to import from another media (or wish to manually insert your data) then simply run `./scripts/make-db.py` and this will create a correctly formatted sqlite database and a config file
    - B. If you do wish to import existing records then I have a Python script under `./scripts/csv2sqlite.py` that you can modify to your purposes. Note this script will create the database for itself.
5. Following this run `scripts/move-configs.py` to move the database and configs to `$HOME/config/cave-logger`.
6. Finally run `make install`.
7. You will now (provided you have a GOPATH set up) be able to run the application by running `cave-logger` in your shell. 

## Help

### Keybindings

| Key | Function |
|:---:|:--------:|
| <kbd>q</kbd> | quit |
| <kbd>n</kbd> | new |
| <kbd>u</kbd> | update |
| <kbd>d</kbd> | delete |
| <kbd>j</kbd> | down |
| <kbd>k</kbd> | up |
| <kbd>g</kbd> | end |
| <kbd>G</kbd> | home |
| <kbd>Tab</kbd> | see menu |
| <kbd>Enter</kbd> | inspect record |

### Menu
In the Menu the <kbd>Tab</kbd> key will select the highlighted item, and hitting <kbd>Tab</kbd> again will navigate to the Menu.

#!/usr/bin/env python3
import json, pathlib, sys

cfg = pathlib.Path(__file__).with_name('config.json')
print(json.load(cfg.open())['root'])
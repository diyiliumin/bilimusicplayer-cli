#!/usr/bin/env python3
import pathlib, sys

root = pathlib.Path(sys.stdin.read().strip())
cid = sys.argv[1]                      # 283912456
d = root / cid
for f in d.glob('*-*.m4s'):            # 283912456_nb2-1-30016.m4s ...
    print(f.resolve())
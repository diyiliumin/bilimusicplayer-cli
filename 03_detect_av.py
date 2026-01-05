#!/usr/bin/env python3
import subprocess, pathlib, sys

f = pathlib.Path(sys.stdin.read().strip())
cmd = ['sh','-c','tail -c +10 "$0" | ffplay -nodisp -autoexit - 2>&1 | grep -m1 Stream',str(f)]
out = subprocess.check_output(cmd,text=True)
print('Audio' if ('mp4a' in out or 'eac3' in out) else 'Video')

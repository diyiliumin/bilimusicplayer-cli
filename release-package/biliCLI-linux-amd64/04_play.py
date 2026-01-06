#!/usr/bin/env python3
import subprocess, sys, os

f = sys.stdin.read().strip()
# 用 mpv 播放，你也可以换成 ffplay
subprocess.run(['tail','-c','+10',f],stdout=subprocess.Popen(['mpv','-'],stdin=subprocess.PIPE).stdin)

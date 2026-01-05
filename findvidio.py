#!/usr/bin/env python3
# bili_tree.py
import json, pathlib, sys, os
from collections import defaultdict

root = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path("/mnt/c/Users/atri_/Videos/bilibili")

groups = defaultdict(list)

# 同时扫两种真实文件名
for name in ('.videoInfo', 'videoInfo.json'):
    for f in root.rglob(name):
        try:
            item = json.load(f.open(encoding='utf-8'))
            item['cid'] = f.parent.name          # ← 把文件夹名（cid）带回来
            groups[item.get('groupTitle', '🚫无合集')].append(item)
        except Exception as e:
            print('skip', f, e, file=sys.stderr)

# 打印
for g_title, items in sorted(groups.items()):
    print(f"{g_title} ({len(items)}P)")
    for v in sorted(items, key=lambda x: int(x.get('p', 0))):
        m, s = divmod(v['duration'], 60)
        size = f"{v['loadedSize']:,d}"
        print(f"  [{v.get('p','?')}] {v['title']}  {m}:{s:02d}  {size}B  #{v['cid']}")

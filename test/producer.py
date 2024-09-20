import redis
import json
import time
from random import randint

"""
r = redis.Redis(host='localhost', port=6379, db=0)
# randomness
for _ in range(5):
    url = f"https://example.com/{randint(1000, 9999)}"
    r.xadd("screenshot-urls", {"url": url})

    time.sleep(0.2)
"""
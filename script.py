#!/usr/bin/env python
import argparse
import redis


def connect_redis(url):
    conn = redis.from_url(url)
    conn.ping()
    return conn


def migrate_redis(source, destination):
    src = connect_redis(source)
    dst = connect_redis(destination)
    for key in src.keys('*'):
        ttl = src.ttl(key)
        # we handle TTL command returning -1 (no expire) or -2 (no key)
        if ttl < 0:
            ttl = 0
        print("Dumping key: %s" % key)
        value = src.dump(key)
        print("Restoring key: %s" % key)
        try:
            dst.restore(key, ttl * 1000, value, replace=True)
        except redis.exceptions.ResponseError as e:
            print("Exception: %s" % str(e))
            print("Failed to restore key: %s" % key)
            pass
    return


def run():
    parser = argparse.ArgumentParser()
    parser.add_argument('source',type=str)
    parser.add_argument('target',type=str)
    options = parser.parse_args()
    migrate_redis(options.source, options.target)

if __name__ == '__main__':
    run()
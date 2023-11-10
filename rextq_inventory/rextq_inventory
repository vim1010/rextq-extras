#!/usr/bin/env python
import requests
from requests.auth import HTTPBasicAuth
import traceback
import json
import os
import sys
import argparse


BASE_URL = os.environ["REX_BASE_URL"]
USER = os.environ["REX_USER"]
PASSWD = os.environ["REX_PASS"]
PROJECT_ID = os.environ["REX_PROJECT_ID"]
auth = HTTPBasicAuth(USER, PASSWD)


def inventory_factory():
    inventory = {
        "_meta": {
            "hostvars": {},
        },
        "all": {
            "hosts": [],
            "vars": {},
            "children": [
                "ungrouped",
            ],
        },
        "ungrouped": {
            "hosts": [],
            "vars": {},
        },
    }
    return inventory


def rpc_stream(fn, payload=None):
    url = f"{BASE_URL}/rpc/{fn}"
    content = []
    res = requests.post(url, json=payload, auth=auth, stream=True)
    content = [str(x) for x in res.iter_content(1024, decode_unicode=True)]
    content = "".join(content)
    d = None
    try:
        d = json.loads(content)
    except Exception as e:
        print(content)
        raise
    return d


def rpc(fn, payload=None):
    url = f"{BASE_URL}/rpc/{fn}"
    res = requests.post(url, json=payload, auth=auth)
    return res.json()


def get_inventory():
    locked_hosts = []
    inventory = inventory_factory()
    g = rpc_stream("get_host_groups", payload=dict(project_id=PROJECT_ID))
    for x in g:
        host_group_id = x["host_group_id"]
        host_group_name = x["host_group_name"]
        if host_group_name == "all":
            print("cannot have host group named 'all', skipping ...", file=sys.stderr)
            continue
        group_vars = x["data"]
        if group_vars is None:
            group_vars = {}
        hosts = []
        opts = dict(host_group_id=host_group_id, project_id=PROJECT_ID)
        res = rpc_stream("get_host_group_ips", payload=opts)
        for h in res:
            host_ip = h["host_ip"]
            if host_ip is None:
                host_id = h["host_id"]
                print(f"host_id [{host_id}] host_ip is None", file=sys.stderr)
                continue
            if h["host_locked"] == True:
                locked_hosts.append(dict(host_group_id=host_group_id, host=host_ip))
            else:
                hosts.append(host_ip)
        inventory[host_group_name] = dict(
            hosts=hosts,
            vars=group_vars,
        )
    d = dict(
        inventory=inventory,
        locked_hosts=locked_hosts,
    )
    return d


def main():
    parser = argparse.ArgumentParser(prog="rextq_inventory")
    parser.add_argument("--list", action="store_true", help="list hosts")
    parser.add_argument("--host", help="stub")
    args = parser.parse_args()
    if args.host:
        return "{}"
    res = get_inventory()
    inventory = json.dumps(res["inventory"])
    print(inventory)


if __name__ == "__main__":

    main()

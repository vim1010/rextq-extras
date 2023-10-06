#!/usr/bin/python
from ansible.module_utils.basic import AnsibleModule
import requests
from requests.auth import HTTPBasicAuth
import os


BASE_URL = os.getenv("REX_BASE_URL", "http://localhost/v1")


def rpc(module):
    params = module.params
    rpc_name = params.get("rpc")
    user = params.get("user")
    passwd = params.get("passwd")
    auth = HTTPBasicAuth(user, passwd)
    payload = module.params.get("payload")
    fn = f"{BASE_URL}/rpc/{rpc_name}"
    res = requests.post(fn, json=payload, auth=auth)
    return res


def main():
    module_args = dict(
        rpc=dict(type="str", required=True),
        user=dict(type="str", required=True),
        passwd=dict(type="str", required=True, no_log=True),
        payload=dict(type="dict"),
    )
    module = AnsibleModule(
        argument_spec=module_args,
        supports_check_mode=True,
    )
    res = rpc(module)
    if res.status_code != 200:
        module.fail_json(msg="request failed", **res.json())
    module.exit_json(changed=False, meta=res.json())


if __name__ == "__main__":
    main()

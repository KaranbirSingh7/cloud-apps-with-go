"""
1. read files: ['.env', '.env-test']
2. run 'az' cli to get resource group name
3. update both files with new subscriptionID and resource groupName
"""
from pathlib import Path
import os
from typing import List
import subprocess as sb
import json

ENV_FILES=[".env-test"]
ROOT_DIR=Path(__file__).parent.parent # this file parent

def az_get_resource_group() -> str:
    res = ""
    cmd = sb.run(args=['az', 'group', 'list'], capture_output=True, universal_newlines=True, check=True)
    cmd_json = json.loads(cmd.stdout)
    res = cmd_json[0].get('name', 'placeholder_rg-HERE')
    return res


def az_get_subscription_id() -> str:
    res = ""
    cmd = sb.run(args=['az', 'account', 'show'], capture_output=True, universal_newlines=True, check=True)
    cmd_json = json.loads(cmd.stdout)
    res = cmd_json.get('id', 'xxx00x0x')
    return res

def read_file(filename:str) -> List[str]:
    res = []
    with open(filename, 'r') as f:
        res = f.readlines()
    return res

def update(sub_id:str, rg_name: str):
    for file in ENV_FILES:
        file_abs_path = Path.joinpath(ROOT_DIR, file)
        is_exists = os.stat(file_abs_path)
        if is_exists:
            print(f"file '{file_abs_path}' exists.")
            contents = read_file(file_abs_path)
            print(contents)



if __name__ == "__main__":
    subscription_id = az_get_subscription_id()
    resource_group = az_get_resource_group()

    output=f"""
AZURE_SUBSCRIPTION_ID={subscription_id}
AZURE_RESOURCE_GROUP_NAME={resource_group}
    """
    print(output)

    # TODO: update script to auto update our ENV files. for now we just utilize output part of it

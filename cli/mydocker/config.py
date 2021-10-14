import os
import json
import click

CONFIG_FILE = os.path.join(os.path.expanduser("~"), ".mydocker.json")
if not os.path.exists(CONFIG_FILE):
    print("Config file will be placed at: %s" % CONFIG_FILE)


class Defaults:
    API_URL = "http://prod.eba-c2qbpgyt.us-east-1.elasticbeanstalk.com"
    LOCAL_API_URL = "http://localhost:5000"
    TOKEN = None
    VERSION = "v1"


def getConfig(key=None) -> dict:
    try:
        with open(CONFIG_FILE, "r") as f:
            contents = json.loads(f.read())
    except FileNotFoundError:
        click.echo("Config file not found. Creating new one...")
        contents = {"token": Defaults.TOKEN, "api_url": Defaults.API_URL}
        with open(CONFIG_FILE, "w") as f:
            f.write(json.dumps(contents))
    if not key:
        return contents
    return contents.get(key, None)


def saveConfig(values: dict):
    contents = getConfig()

    for key in values:
        contents[key] = values[key]

    with open(CONFIG_FILE, "w") as f:
        f.write(json.dumps(contents))

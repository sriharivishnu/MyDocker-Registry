import json
import click
from constants import CONFIG_FILE, API_URL
import requests


def readAndSaveToken(token):
    try:
        with open(CONFIG_FILE, "r") as f:
            contents = json.loads(f.read())
    except Exception:
        contents = {}
    contents[API_URL] = token

    with open(CONFIG_FILE, "w") as f:
        f.write(json.dumps(contents))


def getToken():
    try:
        with open(CONFIG_FILE, "r") as f:
            contents = json.loads(f.read())
            creds = contents[API_URL]
    except Exception:
        raise click.UsageError(
            "Could not find login credentials. Please obtain credentials with the signup or login commands"
        )
    return creds


def doPost(endpoint, payload, withAuth=False):
    if withAuth:
        headers = {"Authorization": "Bearer " + getToken()}
    else:
        headers = {}

    try:
        resp = requests.post(url=API_URL + endpoint, json=payload, headers=headers)
        json_response = resp.json()
        if resp.status_code != 200:
            if "error" in json_response:
                raise Exception(json_response["error"])
            raise Exception("Unknown error occurred while calling API")
    except Exception as e:
        raise click.ClickException(str(e))
    return json_response


def doGet(endpoint, withAuth=False):
    if withAuth:
        headers = {"Authorization": "Bearer " + getToken()}
    else:
        headers = {}

    try:
        resp = requests.get(url=API_URL + endpoint, headers=headers)
        json_response = resp.json()
        if resp.status_code != 200:
            if "error" in json_response:
                raise Exception(json_response["error"])
            raise Exception("Unknown error occurred while calling API")
    except Exception as e:
        raise click.ClickException(str(e))
    return json_response

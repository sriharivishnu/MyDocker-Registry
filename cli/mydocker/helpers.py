import json
import jwt
import click
import requests
import tarfile
import shutil
import os

try:
    from .config import getConfig, saveConfig
except ImportError:
    from config import getConfig, saveConfig


class Token:
    def __init__(self):
        self._cached_token = None

    def readAndSaveToken(self, token: str) -> None:
        saveConfig({"token": token})

    @property
    def token(self) -> str:
        if self._cached_token is not None:
            return self._cached_token

        try:
            creds = getConfig(key="token")
            if creds is None:
                raise Exception("invalid token")
        except Exception:
            raise click.UsageError(
                "Could not find a login token. Please obtain a token with the signup or login command"
            )
        self._cached_token = creds
        return creds

    @property
    def user(self) -> dict:
        return jwt.decode(self.token, options={"verify_signature": False})


def _readResponse(resp):
    try:
        json_response = resp.json()
        if resp.status_code != 200:
            if "error" in json_response:
                raise click.ClickException(json_response["error"])
            raise click.ClickException("Unknown error occurred while calling API")
        return json_response
    except click.ClickException as c:
        raise c
    except Exception as e:
        raise click.ClickException("Unknown error occurred while calling API")


def doPost(endpoint: str, payload: dict, token: Token = None) -> dict:
    if token:
        headers = {"Authorization": "Bearer " + token.token}
    else:
        headers = {}

    try:
        resp = requests.post(
            url=getConfig("api_url") + endpoint, json=payload, headers=headers
        )
        json_response = _readResponse(resp)
    except Exception as e:
        raise click.ClickException(str(e))
    return json_response


def doGet(endpoint: str, token: Token = None) -> dict:
    if token:
        headers = {"Authorization": "Bearer " + token.token}
    else:
        headers = {}
    try:
        resp = requests.get(url=getConfig("api_url") + endpoint, headers=headers)
        json_response = _readResponse(resp)
    except Exception as e:
        raise click.ClickException(str(e))
    return json_response


class Image:
    def __init__(self, name: str):
        self.name = name
        self.user = None
        self.repository = None
        self.tag = "latest"

    def parse(self):
        try:
            parts = self.name.split("/")
            self.user = parts[0]
            if ":" in parts[1]:
                self.repository, self.tag = parts[1].split(":")
            else:
                self.repository = parts[1]

        except Exception:
            raise click.ClickException(
                "Unknown format for image. Please make sure image is in the format <username>/<repository>:tag"
            )
        return self


def zip_tar(tar_file_name):
    TEMP_DIR = ".tmp"
    if not os.path.exists(TEMP_DIR):
        os.mkdir(TEMP_DIR)

    with tarfile.open(tar_file_name) as f:
        f.extractall(TEMP_DIR)

    with tarfile.open(tar_file_name + ".gz", "w:gz") as f:
        f.add(TEMP_DIR, ".")

    shutil.rmtree(TEMP_DIR)
    os.remove(tar_file_name)

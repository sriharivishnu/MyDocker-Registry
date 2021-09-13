import click
import requests
import docker
import getpass
import json
from helpers import readAndSaveToken, getToken, doPost, doGet
from constants import API_URL, CONFIG_FILE


@click.group()
def main():
    """
    CLI application that provides methods for testing the custom docker repository
    """
    pass


@main.command()
def signup():
    """Sign up to push images to the repository"""
    username = input("Username: ")
    password = getpass.getpass()
    if len(username) <= 5 or len(password) <= 6:
        raise click.UsageError(
            message="Invalid username or password. Please make sure username has at least 5 characters, and password has at least 6 characters."
        )

    json_response = doPost("/login", {"username": username, "password": password})

    readAndSaveToken(json_response["token"])
    click.echo(json_response["message"])


@main.command()
def login():
    """Log in to push images to the repository"""
    username = input("Username: ")
    password = getpass.getpass()

    json_response = doPost("/signin", {"username": username, "password": password})
    readAndSaveToken(json_response["token"])
    click.echo(json_response["message"])


@main.command()
@click.argument("image")
def push(image):
    """Pushes image to the repository"""
    doGet("/", withAuth=True)

    client = docker.from_env()
    imageObj = client.images.get(image)
    print(imageObj.attrs)
    # tar_file_name = image.replace(":", "-").replace("/", "_") + ".tar"
    # f = open(tar_file_name, "wb")
    # for chunk in imageObj.save():
    #     f.write(chunk)
    # f.close()


@main.command()
@click.argument("image")
def pull(image):
    """Pulls an image from the repository."""
    print(image)


if __name__ == "__main__":
    main()

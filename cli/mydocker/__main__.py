import click
import docker
import getpass
from tqdm import tqdm
import requests
import os
from prettytable import PrettyTable

try:
    from .helpers import doPost, doGet, Image, Token, zip_tar
except ImportError:
    from helpers import doPost, doGet, Image, Token, zip_tar


@click.group()
def main():
    """
    CLI application that provides methods for testing the custom docker repository
    """
    pass


@main.command()
def signup():
    """Sign up to push images to a repository"""
    username = input("Username: ")
    password = getpass.getpass()
    if len(username) <= 5 or len(password) <= 6:
        raise click.UsageError(
            message="Invalid username or password. Please make sure username has at least 5 characters, and password has at least 6 characters."
        )

    json_response = doPost("/auth/signup", {"username": username, "password": password})

    Token().readAndSaveToken(json_response["token"])
    click.echo(json_response["message"])


@main.command()
def login():
    """Log in to push images to a repository"""
    username = input("Username: ")
    password = getpass.getpass()

    json_response = doPost("/auth/login", {"username": username, "password": password})
    Token().readAndSaveToken(json_response["token"])
    click.echo(json_response["message"])


@main.command()
@click.argument("image")
@click.option("--description", required=False, default="")
def push(image, description):
    """Pushes image to the repository. Expects image in format <username>/<repository>:tag"""
    imageDetails = Image(image).parse()

    client = docker.from_env()
    try:
        dockerImage = client.images.get(image)
    except Exception as e:
        raise click.ClickException(
            "Could not find image locally: %s. Please run docker images to check which images are on your system."
            % image
        )

    tar_file_name = image.replace(":", "-").replace("/", "_") + ".tar"
    with open(tar_file_name, "wb") as f:
        for chunk in dockerImage.save(named=True):
            f.write(chunk)

    zip_tar(tar_file_name)
    tar_file_name += ".gz"

    token = Token()
    upload_url = doGet(
        f"/users/{imageDetails.user}/repositories/{imageDetails.repository}/images/{imageDetails.tag}/upload_url",
        token=token,
    )["upload_url"]

    with open(tar_file_name, "rb") as tar_file:
        try:
            r = requests.put(
                upload_url,
                data=tar_file.read(),
                headers={
                    "Content-Type": "application/gzip",
                    "Content-Transfer-Encoding": "application/gzip",
                },
            )
            r.raise_for_status()
        except Exception as e:
            os.remove(tar_file_name)
            raise click.ClickException(str(e))

    click.echo(
        doPost(
            f"/users/{imageDetails.user}/repositories/{imageDetails.repository}/images",
            payload={"tag": imageDetails.tag, "description": description},
            token=token,
        )["message"]
    )

    os.remove(tar_file_name)


@main.command()
@click.argument("image")
def pull(image):
    """Pulls an image from the repository. Must be in form of <username>/<repository>[:tag]"""
    imageDetails = Image(image).parse()

    imageResponse = doGet(
        f"/users/{imageDetails.user}/repositories/{imageDetails.repository}/images/{imageDetails.tag}"
    )

    download_url = imageResponse["download_url"]
    response = requests.get(download_url, stream=True)
    print(response)
    if response.status_code != 200:
        raise click.ClickException("Could not pull image from repository")

    tar_file_name = imageDetails.name.replace(":", "-").replace("/", "_") + ".tar.gz"

    with open(tar_file_name, "wb") as tarFile:
        for data in tqdm(response.iter_content()):
            tarFile.write(data)

    client = docker.from_env()
    images = client.images.load(open(tar_file_name, "rb").read())

    for image in images:
        click.echo("Created image: %s" % image.id_attribute)

    os.remove(tar_file_name)


@main.command()
@click.argument("repository")
@click.option("--description", required=False, default="")
def create(repository, description):
    """Creates a repository. Please give repository in the format <username>/<repository>"""
    repo = Image(repository).parse()
    doPost(
        f"/users/{repo.user}/repositories",
        {"name": repo.repository, "description": description},
        token=Token(),
    )
    click.echo("Successfully created repository: %s" % repo.name)


@main.command()
@click.argument("user")
def repositories(user):
    """Retrieves the repositories for a given user"""
    response = doGet(f"/users/{user}/repositories")
    repositories = response["repositories"]
    table = PrettyTable()
    table.field_names = ["Name", "Description", "Created"]
    for x in repositories:
        table.add_row([x.get("name"), x.get("description"), x.get("created_at")])

    print(table)


@main.command()
@click.argument("user")
@click.argument("repository")
def images(user, repository):
    """Retrieves the images for a given repository"""
    response = doGet(f"/users/{user}/repositories/{repository}/images")
    images = response["images"]
    table = PrettyTable()
    table.field_names = ["Name", "Description", "Created"]
    for x in images:
        table.add_row(
            [x.get("tag", "<no-tag>"), x.get("description"), x.get("created_at")]
        )

    print(table)


@main.command()
@click.argument("query")
def search(query):
    pass


if __name__ == "__main__":
    main()

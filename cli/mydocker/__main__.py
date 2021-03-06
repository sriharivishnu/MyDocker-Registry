import click
import docker
import getpass
from tqdm import tqdm
import requests
import os
from prettytable import PrettyTable

try:
    from .helpers import doPost, doGet, Image, Token, zip_tar
    from .config import saveConfig, CONFIG_FILE, Defaults
except ImportError:
    from helpers import doPost, doGet, Image, Token, zip_tar
    from config import saveConfig, CONFIG_FILE, Defaults


@click.group()
def main():
    """
    CLI application that provides methods for testing the custom docker repository
    """
    pass


@main.command()
@click.option(
    "--host",
    required=False,
    default=None,
    help="Saves a custom API URL to use for requests",
)
@click.option(
    "--local",
    is_flag=True,
    required=False,
    help="Sets the API URL http://localhost:5000",
)
@click.option(
    "--reset",
    is_flag=True,
    required=False,
    help="Resets configuration to default settings",
)
@click.option(
    "--curuser",
    is_flag=True,
    required=False,
    help="Outputs the current logged in user",
)
def config(host, local, reset, curuser):
    """Set API URL, reset settings, or get current user"""
    if host:
        saveConfig({"api_url": host})
        click.echo("Successfully set new API URL to %s" % host)
    elif local:
        saveConfig({"api_url": Defaults.LOCAL_API_URL, "token": None})
        click.echo("Successfully set new API URL to %s" % Defaults.LOCAL_API_URL)
    elif reset:
        os.remove(CONFIG_FILE)
        click.echo("Successfully reset config to default settings")
    elif curuser:
        token = Token()
        click.echo("Username: %s" % token.user["username"])
    else:
        click.echo("No arguments were supplied. Config was not changed")


@main.command()
def signup():
    """Sign up to push images to a repository"""
    username = input("Username: ")
    password = getpass.getpass()

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
    """Pushes an image to the repository. Expects image in format <username>/<repository>:tag"""
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

    click.echo("Compressing Image...")
    zip_tar(tar_file_name)
    tar_file_name += ".gz"

    token = Token()
    try:
        create_image_response = doPost(
            f"/users/{imageDetails.user}/repositories/{imageDetails.repository}/images",
            payload={"tag": imageDetails.tag, "description": description},
            token=token,
        )
        upload_url = create_image_response["upload_url"]
    except click.ClickException as e:
        os.remove(tar_file_name)
        raise e

    click.echo("Uploading to repository...")
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

    os.remove(tar_file_name)

    click.echo(create_image_response["message"])


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

    if response.status_code != 200:
        raise click.ClickException("Failed to pull image from repository")

    tar_file_name = imageDetails.name.replace(":", "-").replace("/", "_") + ".tar.gz"

    file_size = response.headers.get("Content-length", 0)
    block_size = 1024 * 1024  # 1 MB
    with open(tar_file_name, "wb") as tarFile:
        progress_bar = tqdm(unit="iB", unit_scale=True, total=int(file_size))
        for data in response.iter_content(block_size):
            if data:
                progress_bar.update(len(data))
                tarFile.write(data)
        progress_bar.close()

    client = docker.from_env()
    images = client.images.load(open(tar_file_name, "rb").read())

    for image in images:
        click.echo(
            "\nSuccessfully pulled image: %s" % imageDetails.name
            if len(image.tags) > 0
            else image.id,
        )

    os.remove(tar_file_name)


@main.command()
@click.argument("repository")
@click.option("--description", required=False, default="")
def create(repository, description):
    """Creates a repository"""
    token = Token()
    response = doPost(
        f"""/users/{token.user["username"]}/repositories""",
        {"name": repository, "description": description},
        token=Token(),
    )
    click.echo(response["message"])


@main.command()
@click.argument("user", required=False)
def repositories(user):
    """Retrieves the repositories for a given user. Defaults to logged in user if present"""
    if not user:
        user = Token().user["username"]
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
@click.option("--offset", default=0)
def search(query, offset):
    """Search for image repository"""
    response = doGet(f"/repositories/search?query={query}&offset={offset}")
    repositories = response["results"]
    table = PrettyTable()
    table.field_names = ["User", "Repository", "Description", "# of Images", "Created"]
    for x in repositories:
        table.add_row(
            [
                x.get("username"),
                x.get("repo_name"),
                x.get("description"),
                x.get("num_tags"),
                x.get("created_at"),
            ]
        )

    print(table)


if __name__ == "__main__":
    main()

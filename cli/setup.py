from setuptools import setup, find_packages
from io import open
from os import path

import pathlib

DIR = pathlib.Path(__file__).parent

with open(path.join(DIR, "requirements.txt"), encoding="utf-8") as f:
    all_reqs = f.read().split("\n")

install_requires = [
    x.strip()
    for x in all_reqs
    if ("git+" not in x) and (not x.startswith("#")) and (not x.startswith("-"))
]
dependency_links = [x.strip().replace("git+", "") for x in all_reqs if "git+" not in x]
setup(
    name="mydocker",
    description="A custom docker tool that can push and pull images!",
    version="1.0.0",
    packages=find_packages(),  # list of all packages
    install_requires=install_requires,
    python_requires=">=3.2",
    entry_points="""
        [console_scripts]
        mydocker=mydocker.__main__:main
    """,
    author="Srihari Vishnu",
    license="MIT",
    dependency_links=dependency_links,
    author_email="srihari.vishnu@gmail.com",
)

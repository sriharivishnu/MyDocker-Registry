from setuptools import setup

setup(
    name="mydocker",
    version="0.1.0",
    packages=["mydocker"],
    entry_points={"console_scripts": ["mydocker = mydocker.__main__:main"]},
)

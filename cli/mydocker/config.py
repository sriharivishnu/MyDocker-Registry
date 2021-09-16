import os

API_URL = "http://localhost:5000/v1"
CONFIG_FILE = os.path.join(os.path.expanduser("~"), ".mydocker.json")
if not os.path.exists(CONFIG_FILE):
    print("Config file will be placed at: %s" % CONFIG_FILE)

import os
import toml


config = toml.load("pwgo.toml")
app_name = config["app"]["name"]

app_path = config["app"]["target"]
app_conents = os.path.join(app_path, "Contents")
# exec_name = config["execname"]


def check_existing(config):
    for category in config:
        print(category, type(category))
        if category == "app":
            continue
        print(category)

#!/usr/local/bin/python3

import toml
import os
import sys
import shutil
import subprocess


def create_dirs(path: str, base: str = ""):
    if not os.path.isdir(base):
        try:
            print(f"oh oh, {base} does not exist. attempting to create...")
            os.mkdir(base)
            print("done... close call!")
        except Exception as ex:
            print(ex)
            sys.exit(1)

    for directory in path.strip("/").split("/"):
        base = os.path.join(base, directory)
        if os.path.isdir(base):
            # print(f"{base} already exists")
            continue
        try:
            print(f"attempting to create {base}")
            os.mkdir(base)
        except Exception as ex:
            print(ex)


def main():
    config = toml.load("pwgo.toml")
    origin = config["app"]["source"]
    target = config["app"]["target"]
    app_root = config["app"]["name"] + ".app/Contents/"
    execname = config["app"]["execname"]
    create_dirs(app_root, target)

    # build
    # run_string = f"go build {os.path.join('.', origin)}"
    subprocess.run(
        ["go", "build", os.path.join(".", origin)],
        check=True,
        stdout=subprocess.PIPE,
        universal_newlines=True,
    )
    if execname in os.listdir(origin):
        os.remove(os.path.join(origin, execname))
    shutil.move(os.path.join(".", execname), os.path.join(".", origin))
    for entry in config:
        if entry == "app":
            continue
        print(f"{entry.upper()}")
        entry_data = config.get(entry)
        source = os.path.join(origin, entry_data["source"])
        print(source)
        files = entry_data["files"]
        for file in files:
            destination = ""
            destination = os.path.join(app_root, entry_data["target"])
            create_dirs(destination, target)
            if file not in os.listdir(source):
                print(file)
                print("file missing. exiting...")
                sys.exit(1)
            print(
                f"copying {file} from {source} to {os.path.join(target, destination)}"
            )
            shutil.copy(os.path.join(source, file), os.path.join(target, destination))


if __name__ == "__main__":
    main()

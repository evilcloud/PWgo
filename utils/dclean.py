import json
import os
import sys
import time


def load_settings():
    f = open("dclean.json")
    settings = json.load(f)
    f.close()
    return settings


def is_dir(dir: str) -> bool:
    return os.path.isdir(dir)


def is_file(directory, filename: str) -> bool:
    return os.path.isfile(os.path.join(directory, filename))


def confirm_or_exit(message: str = "") -> bool:
    print(message)
    return input("[ENTER] to continue, [n] to quit: ") == ""


def get_all_txt(directory) -> list:
    return [txtfile for txtfile in os.listdir(directory) if txtfile.endswith(".txt")]


def vertical_print(files: list):
    for name in files:
        print(f"\t{name}")


def find_directory(args: list):
    pfiles = []
    directory = ""
    for item in args:
        if item.startswith("@"):
            directory = item[1:]
            continue
        pfiles.append(item)
    return pfiles, directory


def clean_content(content):
    new_content = ""
    for entry in content:
        new_content += entry.strip() + "\n"
    return new_content


def get_from_file(filepath) -> list:
    """ returns the content of the file """
    with open(filepath, "r") as f:
        content = f.readlines()
    return content


def main():
    tick = time.time()
    settings = load_settings()
    directory = settings["directory"]
    pfiles = settings["files"]

    args, pdirectory = find_directory(sys.argv)

    if pdirectory:
        if is_dir(pdirectory):
            directory = pdirectory
        else:
            print(f"{pdirectory} does not exist, turning to default {directory}")
    else:
        print(f"Default directory {directory}")

    if len(args) > 1:
        pfiles = args[1:]

    print("Target files:")
    vertical_print(pfiles)
    files = [exf for exf in pfiles if is_file(directory, exf)]

    if len(files) == len(pfiles):
        print("All files identified")
    else:
        print("Files identified:")
        vertical_print(files)

    for filename in files:
        filepath = os.path.join(directory, filename)
        content = get_from_file(filepath)
        print(f"\nLoaded {len(content)} entries from {filename}")
        if is_file(directory, filename):
            os.rename(filepath, os.path.join(directory, "old_" + filename))
            print(f"Renamed {filename} to {'old_'+filename}")
        new_content = clean_content(content)

        print(f"Writing the clean content into the file")
        with open(filepath, "w") as f:
            f.writelines(new_content)
        content = get_from_file(filepath)
        print(f"Checked {filename} for {len(content)} lines")
    print(f"All done...\nexecution time {time.time() - tick} sec")


if __name__ == "__main__":
    main()

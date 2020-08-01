import os
import re


directory = "data/"
files = ["adjectives.txt", "nouns.txt"]
alphanum = re.compile(r'')

def get_from_file(filepath: str) -> list:
    """ returns the content of the file """
    with open(filepath, "r") as f:
        content = f.readlines()
    return content


def main():
    for file in files:
        filepath = os.path.join(directory, file)
        
        content = get_from_file(filepath)
        print(f"Loaded {len(content)} entries from {file}")
        
        new_content = ""
        for entry in content:
            new_content += entry.strip() + "\n"
        
        os.rename(filepath, os.path.join(directory, "old_" + file))
        with open(os.path.join(directory, file), "w") as f:
            f.writelines(new_content)
        content = get_from_file(filepath)
        print(f'Checked {file} for {len(content)} lines')


if __name__ == "__main__":
    main()
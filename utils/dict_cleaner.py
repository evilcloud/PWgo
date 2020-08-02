import os
import re
import sys




def get_from_file(filepath: str) -> list:
    """ returns the content of the file """
    with open(filepath, "r") as f:
        content = f.readlines()
    return content



def main():
    directory = "data/"
    files = ["adjectives.txt", "nouns.txt"]
    # alphanum = re.compile(r'')

    args = sys.argv
    if len(args) > 1:
        f = args[1]
        if os.path.isdir(f):
            directory = f
    filepath = os.path.join(directory)
    print(f"Target directory {filepath}")
    sys.exit(0)
    
    for file in files:
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
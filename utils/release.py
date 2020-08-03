import toml


def main():
    config = toml.load("pwgo.toml")

    root = config["app"]["root"]

    for entry in config:
        if entry == "app":
            continue
        print(f"{entry.upper()}")
        entry_data = config.get(entry)
        for item in entry_data:
            print(f"  {item}: {entry_data[item]}")


if __name__ == "__main__":
    main()

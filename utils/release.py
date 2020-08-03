import os
import shutil

import configparser

config = configparser.ConfigParser()
config.read('pwgo.ini')
parser = configparser.ConfigParser()
parser.read_file('pwgo.ini')

for item in config:
    print("-- ", item)
    for i in config[item]:
        print(i)
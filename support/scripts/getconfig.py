#!/usr/bin/env python3

import configparser
import sys
import os

def get_ini_value(file_path, section, key):
    config = configparser.ConfigParser()

    if not os.path.exists(file_path):
        return ""

    try:
        config.read(file_path)

        if not config.has_section(section):
            return ""

        return config.get(section, key) if config.has_option(section, key) else ""

    except Exception as e:
        return ""

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("")
        sys.exit(1)

    _, file_path, section, key = sys.argv
    value = get_ini_value(file_path, section, key)
    print(value)

# Removes duplicate lines from a file
# Usage: python remove_duplicates.py <file>

import sys
import json


def remove_duplicates(file_lines):
    """
    Removes duplicate lines from a file
    """
    lines_set = set()
    for lin in file_lines:
        #if json.loads(lin)["output"] == "":
        #    continue
        lines_set.add(lin)
    return lines_set


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python remove_duplicates.py <file>")
        sys.exit(1)
    orig_file = open(sys.argv[1], "r", encoding="utf-8").readlines()
    lines = remove_duplicates(orig_file)
    file = open("clean_" + sys.argv[1], "w", encoding="utf-8")
    for line in lines:
        file.write(line)
    file.close()
    # Print difference
    print(len(orig_file) - len(lines))

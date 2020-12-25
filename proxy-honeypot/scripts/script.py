# coding: utf8

import sys

from mongodb import MongodbClient


def usage(cmd):
    print("Usage:\n{} password".format(cmd))


def password():
    client = MongodbClient()
    # save password to database
    client.save_password_to_db()
    # sort password,urls,evil_ips


if __name__ == '__main__':
    if len(sys.argv) == 2:
        if sys.argv[1] == "password".lower():
            password()
        else:
            usage(sys.argv[0])
    else:
        usage(sys.argv[0])

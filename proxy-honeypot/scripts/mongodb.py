# coding:utf8


from pymongo import MongoClient
from bson.objectid import ObjectId
from config import config

import urlparse


class MongodbClient(object):
    def __init__(self):
        self.mongodb = MongoClient(config.CONST_MONGODB_CONFIG.get('host'),
                                   config.CONST_MONGODB_CONFIG.get('port'),
                                   unicode_decode_error_handler='ignore',
                                   )

        self.db = self.mongodb[config.CONST_MONGODB_CONFIG.get('database')]
        self.db.authenticate(config.CONST_MONGODB_CONFIG.get('username'),
                             config.CONST_MONGODB_CONFIG.get('password'),
                             source=config.CONST_MONGODB_CONFIG.get('database'),
                             )

        self.collection = self.db[config.CONST_MONGODB_CONFIG.get('collection')]
        self.coll_password = self.db["password"]

    def save_password_to_db(self):
        records = self.collection.find({"flag": 0, "request_parameters": {"$ne": {}}},
                                       no_cursor_timeout=True).batch_size(1)
        for record in records:
            url = record.get('url')
            url_parse = urlparse.urlparse(url)
            site = url_parse.netloc
            from_ip = record.get('origin')
            request_body = record.get('request_body')
            request_header = record.get('request_header')
            header = record.get('header')
            body = record.get('body')
            date_start = record.get('date_start')

            self.collection.update({"_id": ObjectId(record.get("_id"))},
                                   {
                                       "$set": {"flag": 1}
                                   },
                                   True, True
                                   )

            request_parameters = record.get('request_parameters')
            keys = request_parameters.keys()
            intersection = get_intersection(keys, config.CONST_KEYWORD)
            if len(intersection) >= 2:
                ret = dict()
                for i in intersection:
                    t = dict()
                    t[i] = record.get('request_parameters').get(i)[0]
                    ret.update(t)

                value = dict(
                    site=site,
                    url=url,
                    from_ip=from_ip,
                    data=ret,
                    request_parameters=request_parameters,
                    request_header=request_header,
                    request_body=request_body,
                    header=header,
                    body=body,
                    date_start=date_start,
                    status=0,
                )
                print("URL: {}, DATA: {}".format(url, ret))

                self.coll_password.update({"site": site, "data": ret}, value, True, True)
                # self.coll_password.insert(value)

        records.close()

    def clean_password(self):
        self.coll_password.remove({})


def get_intersection(a, b):
    """return intersection of two lists"""
    return list(set(a).intersection(b))

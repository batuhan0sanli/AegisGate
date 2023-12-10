from .__base_object import BaseObject


class Service(BaseObject):
    def __init__(self, name, description, mode, path, security, urls):
        self.name = name
        self.description = description
        self.mode = mode
        self.path = path
        self.security = security
        self.urls = urls

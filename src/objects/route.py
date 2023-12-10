from .__base_object import BaseObject


class Route(BaseObject):
    def __init__(self, name, description, path, allow_methods, target, plugins):
        self.name = name
        self.description = description
        self.path = path
        self.allow_methods = allow_methods
        self.target = target
        self.plugins = plugins

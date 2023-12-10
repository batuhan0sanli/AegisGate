from .__base_object import BaseObject


class Shortcut(BaseObject):
    def __init__(self, name, description, path, method, target_path, target_method, plugins):
        self.name = name
        self.description = description
        self.path = path
        self.method = method
        self.target_path = target_path
        self.target_method = target_method
        self.plugins = plugins

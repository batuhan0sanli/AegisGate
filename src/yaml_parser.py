import os

import yaml
from yamlinclude import YamlIncludeConstructor


class Parser:
    loader = yaml.FullLoader

    def __init__(self, path: str):
        self.path = path
        self.__build()

    def __build(self):
        base_dir = os.path.dirname(os.path.abspath(self.path))
        YamlIncludeConstructor.add_to_loader_class(loader_class=self.loader, base_dir=base_dir)

    def parse(self):
        with open(self.path, 'r', encoding='utf-8') as f:
            data = yaml.load(f, Loader=self.loader)
            return data

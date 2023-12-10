from .yaml_parser import Parser
from src.schemas.authentication import AuthenticationSchema
from src.schemas.route import RouteSchema
from src.schemas.shortcut import ShortcutSchema
from src.schemas.service import ServiceSchema
from src.objects.authentication import AuthenticationFactory
from src.objects.route import Route
from src.objects.shortcut import Shortcut
from src.objects.service import Service
from src.enums.service_modes import ServiceModes


class Configurations:
    def __init__(self, path: str):
        self._parser = Parser(path)
        self._data = self._parser.parse()
        self.authentications = dict()
        self.routes = dict()
        self.shortcuts = dict()
        self.services = dict()

        self.parse_authentications()
        self.parse_routes()
        self.parse_shortcuts()
        self.parse_services()

    def parse_authentications(self):
        authentications = self._data['authentications']
        for name, config in authentications.items():
            serialized_data = AuthenticationSchema().load(config)
            auth = AuthenticationFactory(name=name, **serialized_data).get_authentication()
            self.authentications[name] = auth

    def parse_routes(self):
        routes = self._data['routes']
        for name, config in routes.items():
            serialized_data = RouteSchema().load(config)
            route = Route(name=name, **serialized_data)
            self.routes[name] = route

    def parse_shortcuts(self):
        shortcuts = self._data['shortcuts']
        for name, config in shortcuts.items():
            serialized_data = ShortcutSchema().load(config)
            shortcut = Shortcut(name=name, **serialized_data)
            self.shortcuts[name] = shortcut

    def parse_services(self):
        services = self._data['services']
        for name, config in services.items():
            serialized_data = ServiceSchema().load(config)

            if serialized_data['security'] not in self.authentications:
                raise ValueError(f"Authentication {serialized_data['security']} not found")
            serialized_data['security'] = self.authentications[serialized_data['security']]

            urls = self._get_urls(serialized_data['mode'])
            for url in serialized_data['urls']:
                if url not in urls:
                    raise ValueError(f"'{url}' route not found for '{name}' service")
            serialized_data['urls'] = [urls[url] for url in serialized_data['urls']]

            service = Service(name=name, **serialized_data)
            self.services[name] = service

    def _get_urls(self, mode):
        match mode:
            case ServiceModes.proxy:
                return self.routes
            case ServiceModes.shortcut:
                return self.shortcuts
            case _:
                raise ValueError(f"Unknown mode {mode}")

    def check_service(self, service):
        # Check service security
        if service.security not in self.authentications:
            raise ValueError(f"Authentication {service.security} not found")



configurations = Configurations("configurations/configurations.yaml")

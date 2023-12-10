from .__base_object import BaseObject
from src.enums.authentication_type import AuthenticationType


class BaseAuthentication(BaseObject):
    pass


class NoneAuthentication(BaseAuthentication):
    def __init__(self):
        raise NotImplementedError


class APIKeyAuthentication(BaseAuthentication):
    def __init__(self):
        raise NotImplementedError


class BearerAuthentication(BaseAuthentication):
    def __init__(self):
        raise NotImplementedError


class BasicAuthentication(BaseAuthentication):
    type = 'basic'

    def __init__(self, name, username, password):
        self.name = name
        self.username = username
        self.password = password


class DigestAuthentication(BaseAuthentication):
    def __init__(self):
        raise NotImplementedError


class OAuthAuthentication(BaseAuthentication):
    def __init__(self):
        raise NotImplementedError


class AuthenticationFactory:
    def __init__(self, name, type, config):
        self.name = name
        self.type = type
        self.config = config

    def get_authentication_class(self):
        match self.type:
            case AuthenticationType.NONE:
                return NoneAuthentication
            case AuthenticationType.API_KEY:
                return APIKeyAuthentication
            case AuthenticationType.BEARER:
                return BearerAuthentication
            case AuthenticationType.BASIC:
                return BasicAuthentication
            case AuthenticationType.DIGEST:
                return DigestAuthentication
            case AuthenticationType.OAUTH:
                return OAuthAuthentication
            case _:
                raise NotImplementedError

    def get_authentication(self, *args, **kwargs):
        return self.get_authentication_class()(name=self.name, **self.config)

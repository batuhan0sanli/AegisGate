from .__base_enum import BaseEnum


class AuthenticationType(BaseEnum):
    NONE = "none"
    API_KEY = "api_key"
    BEARER = "bearer"
    BASIC = "basic"
    DIGEST = "digest"
    OAUTH = "oauth"

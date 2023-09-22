from .__base_enum import BaseEnum


class Auth(BaseEnum):
    NONE = "none"
    API_KEY = "api_key"
    BEARER = "bearer"
    BASIC = "basic"
    DIGEST = "digest"
    OAUTH = "oauth"

from .__base_schema import BaseSchema
from marshmallow import fields, pre_load, post_load, Schema
from src.enums.authentication_type import AuthenticationType
from src.objects.authentication import BaseAuthentication, BasicAuthentication


class NoneAuthenticationSchema(BaseSchema):
    def __init__(self):
        super().__init__()
        raise NotImplementedError


class BasicAuthenticationSchema(BaseSchema):
    __obj__ = BasicAuthentication
    username = fields.String(required=True)
    password = fields.String(required=True)


class APIKeyAuthenticationSchema(BaseSchema):
    def __init__(self):
        super().__init__()
        raise NotImplementedError


class BearerAuthenticationSchema(BaseSchema):
    def __init__(self):
        super().__init__()
        raise NotImplementedError


class DigestAuthenticationSchema(BaseSchema):
    def __init__(self):
        super().__init__()
        raise NotImplementedError


class OAuthAuthenticationSchema(BaseSchema):
    def __init__(self):
        super().__init__()
        raise NotImplementedError


class AuthenticationSchema(BaseSchema):
    __obj__ = BaseAuthentication
    type = fields.Enum(AuthenticationType, required=True, by_value=True)
    config = fields.Dict(required=True)

    def get_config_schema(self, type: AuthenticationType):
        match type:
            case AuthenticationType.NONE:
                return NoneAuthenticationSchema
            case AuthenticationType.API_KEY:
                return APIKeyAuthenticationSchema
            case AuthenticationType.BEARER:
                return BearerAuthenticationSchema
            case AuthenticationType.BASIC:
                return BasicAuthenticationSchema
            case AuthenticationType.DIGEST:
                return DigestAuthenticationSchema
            case AuthenticationType.OAUTH:
                return OAuthAuthenticationSchema
            case _:
                raise NotImplementedError

    @post_load
    def post_load(self, data, **kwargs):
        config_schema = self.get_config_schema(data['type'])
        serialized_data = config_schema().load(data['config'])
        data['config'] = serialized_data
        return data

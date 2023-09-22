from dataclasses import dataclass, field
from typing import Union, Any

import marshmallow_dataclass
from marshmallow import validates_schema, ValidationError, EXCLUDE

from src.enums import Auth, RequestComponents


class Base:
    class Meta:
        unknown = EXCLUDE


@dataclass
class NoneAuth(Base):
    pass


@dataclass
class APIKeyAuth(Base):
    key: str = field(metadata={'required': True})
    value: str = field(metadata={'required': True})
    add_to: RequestComponents = field(
        metadata={'load_default': RequestComponents.HEADERS, 'required': False, 'by_value': True})


@dataclass
class BearerAuth(Base):
    token: str = field(metadata={'required': True})


@dataclass
class BasicAuth(Base):
    username: str = field(metadata={'required': True})
    password: str = field(metadata={'required': True})


@dataclass
class AuthConfig(Base):
    type: Auth = field(metadata={'load_default': Auth.NONE, 'required': False, 'by_value': True})
    config: Any = field(metadata={'load_default': {}, 'required': False})

    @validates_schema
    def validate_auth(self, data, **kwargs):
        if data['type'] == Auth.NONE:
            data['config'] = _NoneAuthSchema().load(data['config'])
        elif data['type'] == Auth.API_KEY:
            data['config'] = _APIKeyAuthSchema().load(data['config'])
        elif data['type'] == Auth.BEARER:
            data['config'] = _BearerAuthSchema().load(data['config'])
        elif data['type'] == Auth.BASIC:
            data['config'] = _BasicAuthSchema().load(data['config'])
        else:
            raise ValidationError("Invalid auth type")


_NoneAuthSchema = marshmallow_dataclass.class_schema(NoneAuth)
_APIKeyAuthSchema = marshmallow_dataclass.class_schema(APIKeyAuth)
_BearerAuthSchema = marshmallow_dataclass.class_schema(BearerAuth)
_BasicAuthSchema = marshmallow_dataclass.class_schema(BasicAuth)
AuthSchema = marshmallow_dataclass.class_schema(AuthConfig)

AuthConfigFields = Union[NoneAuth, APIKeyAuth, BearerAuth, BasicAuth]

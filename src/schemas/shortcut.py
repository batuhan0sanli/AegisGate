from .__base_schema import BaseSchema
from marshmallow import fields, pre_load, post_load, Schema
from src.objects.shortcut import Shortcut
from src.enums.http_method import HttpMethod


class ShortcutSchema(BaseSchema):
    __obj__ = Shortcut
    description = fields.String(required=True)
    path = fields.String(required=True)
    method = fields.Enum(HttpMethod)
    target_path = fields.Url(required=True)
    target_method = fields.Enum(HttpMethod)
    plugins = fields.List(fields.String(), allow_none=True)

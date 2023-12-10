from .__base_schema import BaseSchema
from marshmallow import fields, pre_load, post_load, Schema
from src.objects.route import Route
from src.enums.http_method import HttpMethod


class RouteSchema(BaseSchema):
    __obj__ = Route
    description = fields.String(required=True)
    path = fields.String(required=True)
    allow_methods = fields.List(fields.Enum(HttpMethod), required=True)
    target = fields.Url(required=True)
    plugins = fields.List(fields.String(), allow_none=True)

from .__base_schema import BaseSchema
from marshmallow import fields
from src.objects.service import Service
from src.enums.service_modes import ServiceModes


class ServiceSchema(BaseSchema):
    __obj__ = Service
    description = fields.String(required=True)
    mode = fields.Enum(ServiceModes, required=True)
    path = fields.String(required=True)
    security = fields.String(required=True)
    urls = fields.List(fields.String(), required=True)

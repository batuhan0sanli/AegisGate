from .__base_enum import BaseEnum


class RequestComponents(BaseEnum):
    HEADERS = "headers"
    QUERY_PARAMS = "query_params"
    BODY = "body"

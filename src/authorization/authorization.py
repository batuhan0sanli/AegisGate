import requests.auth as requests_auth

import src.authorization.utils as utils
from src.authorization.schemas import AuthSchema, AuthConfigFields
from src.enums import Auth, RequestComponents


class Authorization:
    def __init__(self, auth_config: dict):
        serialized_auth = AuthSchema().load(auth_config)

        self.type: Auth = serialized_auth.type
        self.config: AuthConfigFields = serialized_auth.config

    def get_headers_add_to(self) -> dict:
        if self.type == Auth.API_KEY and self.config.add_to == RequestComponents.HEADERS:
            return {self.config.key: self.config.value}
        elif self.type == Auth.BEARER:
            return {"Authorization": f"Bearer {self.config.token}"}
        elif self.type == Auth.BASIC:
            return {"Authorization": f"Basic {utils.get_basic_token(self.config.username, self.config.password)}"}
        else:
            return {}

    def get_query_params_add_to(self) -> dict:
        if self.type == Auth.API_KEY and self.config.add_to == RequestComponents.QUERY_PARAMS:
            return {self.config.key: self.config.value}
        else:
            return {}

    def get_body_add_to(self) -> dict:
        if self.type == Auth.API_KEY and self.config.add_to == RequestComponents.BODY:
            return {self.config.key: self.config.value}
        else:
            return {}

    def get_requests_auth(self) -> requests_auth.AuthBase:
        if self.type == Auth.BASIC:
            return requests_auth.HTTPBasicAuth(self.config.username, self.config.password)
        elif self.type == Auth.DIGEST:
            return requests_auth.HTTPDigestAuth(self.config.username, self.config.password)

import base64


def get_basic_token(username: str, password: str):
    return base64.b64encode(f"{username}:{password}".encode()).decode("ascii")

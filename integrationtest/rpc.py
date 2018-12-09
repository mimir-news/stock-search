# Standard library
from typing import Dict

# 3rd party modules
import requests


def make_request(test_data: Dict, env: Dict[str, str]) -> requests.Response:
    method = test_data["method"]
    if method in ("GET", "DELETE"):
        return _make_request_without_body(test_data, env)
    return _make_request_with_body(test_data, env)


def _make_path(test_data: Dict, env: Dict[str, str]) -> str:
    path = env["baseUrl"] + test_data["path"]
    for key, val in env.items():
        path_key = "${" + key + "}"
        path = path.replace(path_key, val)
    return path


def _make_request_without_body(
    test_data: Dict, env: Dict[str, str]
) -> requests.Response:
    url = _make_path(test_data, env)
    method = test_data["method"]
    headers = _make_headers(env, test_data["withToken"])
    return requests.request(method, url, headers=headers)


def _make_request_with_body(test_data: Dict, env: Dict[str, str]) -> requests.Response:
    if not "body" in test_data:
        return _make_request_without_body(test_data, env)
    url = _make_path(test_data, env)
    method = test_data["method"]
    headers = _make_headers(env, test_data["withToken"])
    body = _resolve_body(env, test_data["body"])
    return requests.request(method, url, json=body, headers=headers)


def _make_headers(env: Dict[str, str], with_token: bool) -> Dict[str, str]:
    base_headers = {"Content-Type": "application/json",
                    "X-Client-ID": env["clientId"]}
    if with_token:
        base_headers["Authorization"] = f'Bearer {env["authToken"]}'
    return base_headers


def _resolve_body(env: Dict[str, str], body: Dict) -> Dict:
    final_body = {}
    for key, val in body.items():
        if isinstance(val, str):
            final_body[key] = _resolve_value(env, val)
        elif isinstance(val, dict):
            final_body[key] = _resolve_body(env, val)
        else:
            final_body[key] = val
    return final_body


def _resolve_value(env: Dict[str, str], val: str) -> str:
    key = val.replace("${", "").replace("}", "")
    if key == val:
        return val
    return env[key]

import json

body = response.json()

if "access_token" in body:
    posting.set_env("TOKEN", body["access_token"])

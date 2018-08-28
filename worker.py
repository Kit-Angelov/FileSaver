import os
import requests
import json

r = requests.post("http://192.168.2.243:8000/rest-auth/login/", data={"username": "admin", "password": "itt0root"})
print(r.text)
data = json.loads(r.text)
token = data["token"]

headers = {"Authorization": "bearer {}".format(token)}

f = open("image.jpg", "rb")
files = {"file": f}

data = {"name": "first.jpg"}

p = requests.post("http://localhost:8080/api/v1/upload/", headers=headers, files=files, data=data)
print(p.text)
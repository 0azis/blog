from time import time
import mysql.connector as mc
import random
import string
import requests
import json
import time

def upload(file: str):
    token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlIjoxMDAwMDAwMDAwMDAwMDAwMCwiaWQiOjIzfQ.2q9WyyBgIGNq3b4ZIGLjRqy9wUZ8p2tV16m1N1gauKA"
    multipart_form_data = {
        'file': open(file, 'rb')
    }
    r = requests.post('http://localhost:8000/api/v1/uploads', files=multipart_form_data, headers={"Authorization": f"Bearer {token}"})
    data = json.loads(r.content)
    return data['filename']

def generate_string(len: int):
    return ''.join(random.choice(string.ascii_lowercase) for _ in range(len))

cnx = mc.connect(
    host="localhost",
    port=3306,
    user="admin",
    password="test123",
    database="blogdb",
)
cur = cnx.cursor()

# create 10 users
def create_users(image_name: str):
    for x in range (0, 10):
        username = generate_string(6)
        name = generate_string(10)
        email = generate_string(random.choice(range(6,10))) + "@gmail.com"
        password = generate_string(10)
        avatar = image_name
        cur.execute(f"insert into users (email, username, password, name, avatar) values ('{email}', '{username}', '{password}', '{name}', '{avatar}')")
        cnx.commit()

def create_posts(preview_name: str):
    for x in range(0, 10):
        userid = random.choice(range(1, 11))
        title = "This is post!"
        preview = preview_name
        content = "<h1>My Post</h1>"
        public = 1
        cur.execute(f"insert into posts (user_id, title, preview, content, public) values ({userid}, '{title}', '{preview}', '{content}', {public})")
        cnx.commit()

def create_comments():
    for x in range(1, 10):
        for y in range(0, 2):
            userid = random.choice(range(1, 11))
            postid = x
            text = "This post is shit!"
            cur.execute(f"insert into comments (post_id, user_id, text) values ({postid}, {userid}, '{text}')")
            cnx.commit()

if __name__ == "__main__":
    avatar = upload('avatar.jpg')
    preview = upload('preview.png')
    print("Photos uploaded!")
    create_users(avatar)
    print("Users created!")
    create_posts(preview)
    print("Posts created!")
    create_comments()
    print("Comments created!")
    cnx.close()
    print("Connection closed!")

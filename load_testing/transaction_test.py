from locust import SequentialTaskSet, HttpUser, constant, task, between
from random import randrange
import names
import time

class Tasks(SequentialTaskSet):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._local_user = None
        self._local_id = None

    def on_start(self):
        # Create User and id
        name = names.get_full_name()
        email = name.replace(" ","") + '@example.com'
        self._local_user = {"name": name, "email": email, "password": "password123"}
        id = f'user{randrange(100000)}'
        self._local_id = '{'+f'{id}'+'}'
        res = self.client.post("/users", json=self._local_user)
        print(f'Created User! response: {res.text}')

    @task
    def post_transaction(self):
        # Post a transaction
        timestamp = int(time.time())
        transaction = {
            "userId": self._local_id, "amount": 3.70, 
            "category": "coffee", "timestamp": timestamp
        }
        res = self.client.post("/transactions", json=transaction)
        print(f'response: {res.text}')

class LoadTest(HttpUser):
    host = 'http://localhost:8080'
    tasks = [Tasks]
    wait_time  = between(1,60)
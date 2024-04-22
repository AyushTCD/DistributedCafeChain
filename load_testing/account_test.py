from locust import SequentialTaskSet, HttpUser, constant, task, between
from random import randrange, choice, choices
import names
import time

class Tasks(SequentialTaskSet):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._local_ids = []
        self._local_accounts = []
        self._local_users = []
        # 2 - 6 per group
        self._group_size = randrange(2,7)

    def on_start(self) -> None:
        
        # Create Users, ids and publich them
        for _ in range(self._group_size):
            time.sleep(0.5)
            name = names.get_full_name()
            email = name.replace(" ","") + '@example.com'
            local_user = {"name": name, "email": email, "password": "password123"}
            number = randrange(100000)
            id = f'user{number}'
            account = f'account{number}'
            local_account = '{'+f'{account}'+'}'
            local_id = '{'+f'{id}'+'}'
            res = self.client.post("/users", json=local_user)
            print(f'Created User! response: {res.text}')
            self._local_users.append(local_user)
            self._local_ids.append(local_id)
            self._local_accounts.append(local_account)
        
        group = {
            'userIds': self._local_ids, 'points': 100,
            "discounts": ["10% off", "20% off"]
        }
        time.sleep(0.5)
        res = self.client.post("/loyalty-accounts", json=group)
        print(res.text)
    
    @task
    def post_transaction(self):
        # Post a transaction
        timestamp = int(time.time())
        id = choice(self._local_ids)
        transaction = {
            "userId": id, "amount": 3.70, 
            "category": "coffee", "timestamp": timestamp
        }
        res = self.client.post("/transactions", json=transaction)
        print(f'response: {res.text}')

class LoadTest(HttpUser):
    host = 'http://localhost:8080'
    tasks = [Tasks]
    wait_time  = between(1,5)


        
        
    
    
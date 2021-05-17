import os

from decouple import config

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

BASE_PATH = config("BASE_PATH", default="/")

DATABASE_URI = (
    f"mysql+pymysql://{config('DB_USER', default='discounts')}:{config('DB_PASSWORD', default='my-password')}@"
    f"{config('DB_HOST', default='localhost:3306')}/{config('DB_NAME', default='discounts')}"
    "?charset=utf8"
)

USERS_EVENTS = config("USERS_EVENTS", default='http://localhost:4566/000000000000/Discounts-UsersCrud.fifo')

PRODUCTS_EVENTS = config("PRODUCTS_EVENTS", default='http://localhost:4566/000000000000/Discounts-ProductsCrud.fifo')

AWS_REGION = config("AWS_REGION", default='sa-east-1')

AWS_ACCESS_KEY = config("AWS_ACCESS_KEY", default='fake_access_key')

AWS_SECRET_KEY = config("AWS_SECRET_KEY", default='fake_secret_key')

AWS_ENDPOINT = config("AWS_ENDPOINT", default='http://localhost:4566')

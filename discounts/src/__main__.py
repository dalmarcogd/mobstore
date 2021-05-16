import logging

from app import Server
from src.app import Consumer

if __name__ == '__main__':
    logging.info("Starting application.")
    Consumer.run()
    Server.run()

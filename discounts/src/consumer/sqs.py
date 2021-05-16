import logging
from typing import Callable

import boto3

from src import settings

_sqs = boto3.client('sqs', region_name=settings.AWS_REGION,
                    aws_access_key_id=settings.AWS_ACCESS_KEY,
                    aws_secret_access_key=settings.AWS_SECRET_KEY,
                    endpoint_url=settings.AWS_ENDPOINT)


def start_pool(queue: str, handler: Callable):
    while True:
        try:
            response = _sqs.receive_message(
                QueueUrl=queue,
                MaxNumberOfMessages=1,
                MessageAttributeNames=[
                    'All'
                ],
                WaitTimeSeconds=10
            )
            if 'Messages' in response:
                try:
                    messages = response['Messages']
                    for message in messages:
                        handler(message)
                except Exception as e:
                    logging.error(f'[sqs] error no message in queue -> {e}')
        except BaseException as exc:
            logging.error(f'[sqs] error no message in queue -> {exc}')

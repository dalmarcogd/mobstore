import logging
from typing import Dict

from src.database import queries
from src.exceptions.exceptions import UnrecognizedEventOperation


def _get_user(message: Dict) -> Dict:
    return {
        'id': message.get('user_id'),
        'first_name': message.get('first_name'),
        'last_name': message.get('last_name'),
        'birth_date': message.get('birth_date'),
    }


def _handle_create_user(user: Dict):
    queries.create_user(user)


def _handle_update_user(user: Dict):
    queries.update_user(user.get('id'), user)


def _handle_delete_user(user: Dict):
    queries.delete_user(user.get('id'))


def handle_users_events(message: Dict):
    event_type: str = message.get('event_type')
    operation: str = message.get('operation')

    if event_type == 'users':
        if operation == 'create':
            user = _get_user(message)
            _handle_create_user(user)
            logging.info(f"user id={user.get('id')} created")
        elif operation == 'update':
            user = _get_user(message)
            _handle_update_user(user)
            logging.info(f"user id={user.get('id')} updated")
        elif operation == 'delete':
            user = _get_user(message)
            _handle_delete_user(user)
            logging.info(f"user id={user.get('id')} deleted")
        else:
            raise UnrecognizedEventOperation(operation)
    else:
        raise UnrecognizedEventOperation(event_type)

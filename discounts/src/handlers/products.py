import logging
from typing import Dict

from src.database import queries
from src.exceptions.exceptions import UnrecognizedEventOperation


def _get_product(message: Dict) -> Dict:
    return {
        'id': message.get('product_id'),
        'title': message.get('title'),
        'description': message.get('description'),
        'price_in_cents': message.get('price_in_cents'),
    }


def _handle_create_product(product: Dict):
    queries.create_product(product)


def _handle_update_product(product: Dict):
    queries.update_product(product.get('id'), product)


def _handle_delete_product(product: Dict):
    queries.delete_product(product.get('id'))


def handle_products_events(message: Dict):
    event_type: str = message.get('event_type')
    operation: str = message.get('operation')

    if event_type == 'products':
        if operation == 'create':
            product = _get_product(message)
            _handle_create_product(product)
            logging.info(f"product id={product.get('id')} created")
        elif operation == 'update':
            product = _get_product(message)
            _handle_update_product(product)
            logging.info(f"product id={product.get('id')} updated")
        elif operation == 'delete':
            product = _get_product(message)
            _handle_delete_product(product)
            logging.info(f"product id={product.get('id')} deleted")
        else:
            raise UnrecognizedEventOperation(operation)
    else:
        raise UnrecognizedEventOperation(event_type)

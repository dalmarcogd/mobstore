from datetime import datetime
from typing import Dict

from sqlalchemy.orm import Query

from src.database.models import User
from src.database.models.products import Product
from src.database.utils import db_session
from src.exceptions.exceptions import (
    UserNotFoundException, ProductNotFoundException,
)


def create_user(user: Dict) -> Dict:
    with db_session() as session:
        new_user = User(**user)
        session.add(new_user)
        session.commit()
        return new_user.to_dict()


def update_user(user_id: str, user: Dict) -> Dict:
    with db_session() as session:
        user_update: User = Query(User, session=session).filter_by(id=user_id).first()
        if not user_update:
            raise UserNotFoundException()
        if 'first_name' in user:
            user_update.first_name = user.get('first_name')
        if 'last_name' in user:
            user_update.last_name = user.get('last_name')
        if 'birth_date' in user:
            user_update.birth_date = user.get('birth_date')
        if 'deleted_at' in user:
            user_update.deleted_at = user.get('deleted_at')

        session.commit()

    return get_user(user_id)


def delete_user(user_id: str) -> [Dict, None]:
    with db_session() as session:
        user_update: User = Query(User, session=session).filter_by(id=user_id).first()
        if not user_update:
            raise UserNotFoundException()

        result = update_user(user_id, {'deleted_at': datetime.utcnow()})
        session.commit()
        return result
    return None


def get_user(user_id: str) -> Dict:
    with db_session() as session:
        user: User = Query(User, session=session).filter_by(id=user_id).one()
        return user.to_dict()


def create_product(product: Dict) -> Dict:
    with db_session() as session:
        new_product = Product(**product)
        session.add(new_product)
        session.commit()
        return new_product.to_dict()


def update_product(product_id: str, product: Dict) -> Dict:
    with db_session() as session:
        product_update: Product = Query(Product, session=session).filter_by(id=product_id).first()
        if not product_update:
            raise ProductNotFoundException()
        if 'title' in product:
            product_update.title = product.get('title')
        if 'description' in product:
            product_update.description = product.get('description')
        if 'price_in_cents' in product:
            product_update.price_in_cents = product.get('price_in_cents')
        if 'deleted_at' in product:
            product_update.deleted_at = product.get('deleted_at')

        session.commit()

    return get_product(product_id)


def delete_product(product_id: str) -> [Dict, None]:
    with db_session() as session:
        product_update: Product = Query(Product, session=session).filter_by(id=product_id).first()
        if not product_update:
            raise ProductNotFoundException()

        result = update_product(product_id, {'deleted_at': datetime.utcnow()})
        session.commit()
        return result
    return None


def get_product(product_id: str) -> Dict:
    with db_session() as session:
        product: Product = Query(Product, session=session).filter_by(id=product_id).one()
        return product.to_dict()

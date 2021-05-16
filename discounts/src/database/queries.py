from typing import Dict

from sqlalchemy.orm import Query

from src.database.models.discounts import User
from src.database.utils import db_session
from src.exceptions.discounts import (
    UserNotFoundException,
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

        result = update_user(user_id, {'deleted_at': True})
        session.commit()
        return result
    return None


def get_user(user_id: str) -> Dict:
    with db_session() as session:
        user: User = Query(User, session=session).filter_by(id=user_id).one()
        return user.to_dict()

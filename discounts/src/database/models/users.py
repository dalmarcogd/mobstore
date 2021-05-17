from typing import Dict

from sqlalchemy import (
    Column,
    String,
    Date,
)

from src.database.models.base import BaseModel


class User(BaseModel):
    __tablename__ = "users"

    first_name = Column(String(100), nullable=False)
    last_name = Column(String(500), nullable=False)
    birth_date = Column(Date, nullable=False)

    def to_dict(self) -> Dict:
        return {
            "id": self.id,
            "created_at": self.created_at,
            "updated_at": self.updated_at,
            "deleted_at": self.deleted_at,
            "title": self.title,
            "description": self.description,
            "price_in_cents": self.price_in_cents,
        }

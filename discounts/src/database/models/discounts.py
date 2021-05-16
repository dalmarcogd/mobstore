from enum import Enum
from typing import Dict

from sqlalchemy import (
    Column,
    Time,
    DECIMAL,
    Integer,
    String,
    ForeignKey, Date,
)
from sqlalchemy.orm import relationship

from src.database.models.base import BaseModel


class User(BaseModel):
    __tablename__ = "users"

    first_name = Column(String, nullable=False)
    last_name = Column(String, nullable=False)
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

class Product(BaseModel):
    __tablename__ = "products"

    title = Column(String, nullable=False)
    description = Column(String, nullable=False)
    price_in_cents = Column(DECIMAL, nullable=False)

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

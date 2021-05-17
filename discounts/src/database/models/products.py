from typing import Dict

from sqlalchemy import (
    Column,
    String, Integer,
)

from src.database.models.base import BaseModel


class Product(BaseModel):
    __tablename__ = "products"

    title = Column(String(100), nullable=False)
    description = Column(String(500), nullable=False)
    price_in_cents = Column(Integer, nullable=False)

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

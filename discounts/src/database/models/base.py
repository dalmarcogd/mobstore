import abc
from datetime import datetime
from typing import Dict

from sqlalchemy import Column, DateTime, String
from sqlalchemy.ext.declarative import declarative_base

DeclarativeModel = declarative_base()


class BaseModel(DeclarativeModel):
    __abstract__ = True
    id = Column(
        String(56),
        unique=True,
        nullable=False,
        primary_key=True,
        index=True,
    )
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)
    updated_at = Column(
        DateTime, nullable=True, default=None, onupdate=datetime.utcnow()
    )
    deleted_at = Column(DateTime, nullable=True, default=None)

    @abc.abstractmethod
    def to_dict(self) -> Dict:
        raise NotImplementedError("Implement method `to_dict`, must be return a Dict.")

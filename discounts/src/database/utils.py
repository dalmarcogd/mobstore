from contextlib import contextmanager

from sqlalchemy import create_engine
from sqlalchemy.orm import scoped_session, sessionmaker, Session

from src.settings import DATABASE_URI

db_engine = create_engine(DATABASE_URI)
db_session_factory = sessionmaker(bind=db_engine)
db_scoped_session = scoped_session(db_session_factory)


@contextmanager
def db_session() -> Session:
    session: Session = db_scoped_session()
    try:
        yield session
        session.commit()
    except Exception as e:
        session.rollback()
        raise e
    finally:
        session.close()

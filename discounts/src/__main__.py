import logging
import sys

from src.app import Server
from src.app import Consumer
from src.database import migration
from src.exceptions.exceptions import UnrecognizedArgs

if __name__ == '__main__':
    logging.root.setLevel(logging.INFO)
    if len(sys.argv) == 1:
        logging.info("Starting application.")
        Consumer.run()
        Server.run()
    elif len(sys.argv) == 2:
        arg = sys.argv[1]
        if arg == 'database_migration':
            migration.upgrade()
        else:
            raise UnrecognizedArgs(f'Unrecognized args (${sys.argv})')
    else:
        raise UnrecognizedArgs(f'Unrecognized args (${sys.argv})')

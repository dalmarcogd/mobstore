from concurrent import futures

import grpc

from src import settings
from src.consumer import sqs
from src.discountsgrpc import discounts_pb2_grpc
from src.handlers.disounts import Discounts
from src.handlers.products import handle_products_events
from src.handlers.users import handle_users_events


class Server:

    @staticmethod
    def run():
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        discounts_pb2_grpc.add_DiscountsServicer_to_server(Discounts(), server)
        server.add_insecure_port('[::]:50051')
        server.start()
        server.wait_for_termination()


class Consumer:

    @staticmethod
    def run():
        ex = futures.ThreadPoolExecutor(max_workers=2)
        ex.submit(sqs.start_pool, settings.PRODUCTS_EVENTS, handle_products_events)
        ex.submit(sqs.start_pool, settings.USERS_EVENTS, handle_users_events)

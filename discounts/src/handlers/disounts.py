from src.discountsgrpc import discounts_pb2_grpc

from src.discountsgrpc.domains_pb2 import DiscountRequest, DiscountResponse


class Discounts(discounts_pb2_grpc.DiscountsServicer):

    def Get(self, request: DiscountRequest, context) -> DiscountResponse:
        return DiscountResponse()

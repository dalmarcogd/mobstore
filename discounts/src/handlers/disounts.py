from datetime import datetime

import grpc
from pymysql import Date

from src.database import queries
from src.discountsgrpc import discounts_pb2_grpc
from src.discountsgrpc.domains_pb2 import DiscountRequest, DiscountResponse


class Discounts(discounts_pb2_grpc.DiscountsServicer):

    def Get(self, request: DiscountRequest, context) -> DiscountResponse:
        user = queries.get_user(request.user_id)
        if not user:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('User not found')
            return DiscountResponse()
        product = queries.get_product(request.product_id)
        if not product:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('Product not found')
            return DiscountResponse()
        percentage = 0
        value_in_cents = product.get("price_in_cents", 0)
        now = datetime.utcnow()
        if user.get('birth_date'):
            birth_date: Date = user.get('birth_date')
            if now.month == birth_date.month and now.day == birth_date.day:
                percentage = 5
        if now.month == 11 and now.day == 25:
            percentage = 10
        if percentage > 10:
            percentage = 10

        if value_in_cents > 0:
            value_in_cents = value_in_cents * (1 - (percentage / 100))
        return DiscountResponse(user_id=request.user_id, product_id=request.product_id, percentage=percentage,
                                value_in_cents=value_in_cents)

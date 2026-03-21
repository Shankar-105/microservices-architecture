import asyncio
import logging
import sys
from pathlib import Path
from typing import Any, cast

import grpc

from app.service.payment_service import PaymentService

ROOT = Path(__file__).resolve().parents[2]
GENERATED_PY = ROOT / "generated-py"
if str(GENERATED_PY) not in sys.path:
    sys.path.insert(0, str(GENERATED_PY))

from bookstore import common_pb2, payment_pb2, payment_pb2_grpc

common_pb2_any = cast(Any, common_pb2)
payment_pb2_any = cast(Any, payment_pb2)

logger = logging.getLogger(__name__)


class PaymentGrpcServicer(payment_pb2_grpc.PaymentServiceServicer):
    def __init__(self, service: PaymentService) -> None:
        self._service = service

    async def Health(self, request, context):
        return common_pb2_any.HealthCheckResponse(status=self._service.health(), service="payment-service-py")

    async def Authorize(self, request, context):
        if not request.order_id:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("order_id is required")
            return payment_pb2_any.AuthorizePaymentResponse()

        decision = self._service.authorize(request.order_id, request.amount_cents)
        return payment_pb2_any.AuthorizePaymentResponse(
            approved=decision.approved,
            transaction_id=decision.transaction_id,
            reason=decision.reason,
        )


async def start_grpc_server(port: int, service: PaymentService) -> None:
    server = grpc.aio.server()
    payment_pb2_grpc.add_PaymentServiceServicer_to_server(PaymentGrpcServicer(service), server)
    listen_addr = f"[::]:{port}"
    server.add_insecure_port(listen_addr)
    await server.start()
    logger.info("payment-service gRPC listening on %s", listen_addr)

    try:
        await server.wait_for_termination()
    except asyncio.CancelledError:
        logger.info("payment-service gRPC shutdown requested")
        await server.stop(grace=5)
        raise

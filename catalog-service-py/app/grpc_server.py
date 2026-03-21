import asyncio
import logging
import sys
from pathlib import Path
from typing import Any, cast

import grpc

from app.service.catalog_service import CatalogService

ROOT = Path(__file__).resolve().parents[2]
GENERATED_PY = ROOT / "generated-py"
if str(GENERATED_PY) not in sys.path:
    sys.path.insert(0, str(GENERATED_PY))

from bookstore import catalog_pb2, catalog_pb2_grpc, common_pb2

catalog_pb2_any = cast(Any, catalog_pb2)
common_pb2_any = cast(Any, common_pb2)

logger = logging.getLogger(__name__)


class CatalogGrpcServicer(catalog_pb2_grpc.CatalogServiceServicer):
    def __init__(self, service: CatalogService) -> None:
        self._service = service

    async def Health(self, request, context):
        return common_pb2_any.HealthCheckResponse(status=self._service.health(), service="catalog-service-py")

    async def GetBook(self, request, context):
        if not request.book_id:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("book_id is required")
            return catalog_pb2_any.GetBookResponse()

        book = self._service.get_book(request.book_id)
        if book is None:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("book not found")
            return catalog_pb2_any.GetBookResponse()

        return catalog_pb2_any.GetBookResponse(
            book_id=book.book_id,
            title=book.title,
            price_cents=book.price_cents,
            available=book.available,
        )


async def start_grpc_server(port: int, service: CatalogService) -> None:
    server = grpc.aio.server()
    catalog_pb2_grpc.add_CatalogServiceServicer_to_server(CatalogGrpcServicer(service), server)
    listen_addr = f"[::]:{port}"
    server.add_insecure_port(listen_addr)
    await server.start()
    logger.info("catalog-service gRPC listening on %s", listen_addr)

    try:
        await server.wait_for_termination()
    except asyncio.CancelledError:
        logger.info("catalog-service gRPC shutdown requested")
        await server.stop(grace=5)
        raise

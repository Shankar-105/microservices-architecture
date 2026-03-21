import asyncio
import logging

logger = logging.getLogger(__name__)


async def start_grpc_server(port: int) -> None:
    # Phase 1: this simulates gRPC server lifecycle while contracts are being finalized.
    logger.info("catalog-service gRPC placeholder listening on :%s", port)
    while True:
        await asyncio.sleep(3600)

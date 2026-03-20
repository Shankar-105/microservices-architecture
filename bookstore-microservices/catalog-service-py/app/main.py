import asyncio
import logging

import uvicorn
from fastapi import FastAPI

from app.config import settings
from app.grpc_server import start_grpc_server
from app.logging_config import configure_logging

configure_logging()
logger = logging.getLogger(__name__)

app = FastAPI(title="catalog-service")


@app.get("/healthz")
async def healthz() -> dict[str, str]:
    return {"service": "catalog-service", "status": "ok"}


async def _serve_http() -> None:
    config = uvicorn.Config(app=app, host="0.0.0.0", port=settings.http_port, log_level="info")
    server = uvicorn.Server(config)
    await server.serve()


async def main() -> None:
    logger.info("starting catalog-service")
    await asyncio.gather(
        start_grpc_server(settings.grpc_port),
        _serve_http(),
    )


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("catalog-service stopped")

import os


def positive_port(name: str, fallback: int) -> int:
    raw = os.getenv(name)
    if raw is None:
        return fallback
    try:
        port = int(raw)
    except ValueError as exc:
        raise RuntimeError(f"{name} must be an integer") from exc
    if port < 1 or port > 65535:
        raise RuntimeError(f"{name} must be between 1 and 65535")
    return port


def required_db_path(name: str, fallback: str) -> str:
    value = os.getenv(name, fallback).strip()
    if not value or "\x00" in value:
        raise RuntimeError(f"{name} must contain a valid database path")
    return value


class Settings:
    grpc_port: int = positive_port("CATALOG_GRPC_PORT", 50053)
    http_port: int = positive_port("CATALOG_HTTP_PORT", 8001)
    db_path: str = required_db_path("CATALOG_DB_PATH", "catalog-service.db")


settings = Settings()

import os


class Settings:
    grpc_port: int = int(os.getenv("CATALOG_GRPC_PORT", "50053"))
    http_port: int = int(os.getenv("CATALOG_HTTP_PORT", "8001"))
    db_path: str = os.getenv("CATALOG_DB_PATH", "catalog-service.db")


settings = Settings()

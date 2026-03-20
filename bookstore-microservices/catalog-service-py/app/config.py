import os


class Settings:
    grpc_port: int = int(os.getenv("CATALOG_GRPC_PORT", "50053"))
    http_port: int = int(os.getenv("CATALOG_HTTP_PORT", "8001"))


settings = Settings()

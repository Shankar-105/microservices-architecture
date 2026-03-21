import os


class Settings:
    grpc_port: int = int(os.getenv("PAYMENT_GRPC_PORT", "50054"))
    http_port: int = int(os.getenv("PAYMENT_HTTP_PORT", "8002"))
    db_path: str = os.getenv("PAYMENT_DB_PATH", "payment-service.db")


settings = Settings()

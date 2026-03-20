import os


class Settings:
    grpc_port: int = int(os.getenv("PAYMENT_GRPC_PORT", "50054"))
    http_port: int = int(os.getenv("PAYMENT_HTTP_PORT", "8002"))


settings = Settings()

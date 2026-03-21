import sqlite3


class PaymentRepository:
    def __init__(self, db_path: str) -> None:
        self._db_path = db_path
        self._init_db()

    def _connect(self) -> sqlite3.Connection:
        return sqlite3.connect(self._db_path)

    def _init_db(self) -> None:
        with self._connect() as conn:
            conn.execute(
                """
                CREATE TABLE IF NOT EXISTS payments (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    order_id TEXT NOT NULL,
                    amount_cents INTEGER NOT NULL,
                    approved INTEGER NOT NULL,
                    transaction_id TEXT NOT NULL,
                    reason TEXT NOT NULL,
                    created_at TEXT DEFAULT CURRENT_TIMESTAMP
                )
                """
            )
            conn.commit()

    def create(self, order_id: str, amount_cents: int, approved: bool, transaction_id: str, reason: str) -> None:
        with self._connect() as conn:
            conn.execute(
                """
                INSERT INTO payments (order_id, amount_cents, approved, transaction_id, reason)
                VALUES (?, ?, ?, ?, ?)
                """,
                (order_id, amount_cents, 1 if approved else 0, transaction_id, reason),
            )
            conn.commit()

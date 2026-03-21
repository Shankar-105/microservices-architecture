import sqlite3
from dataclasses import dataclass


@dataclass
class Book:
    book_id: str
    title: str
    price_cents: int
    available: bool


class CatalogRepository:
    def __init__(self, db_path: str) -> None:
        self._db_path = db_path
        self._init_db()

    def _connect(self) -> sqlite3.Connection:
        return sqlite3.connect(self._db_path)

    def _init_db(self) -> None:
        with self._connect() as conn:
            conn.execute(
                """
                CREATE TABLE IF NOT EXISTS books (
                    book_id TEXT PRIMARY KEY,
                    title TEXT NOT NULL,
                    price_cents INTEGER NOT NULL,
                    available INTEGER NOT NULL
                )
                """
            )

            existing_count = conn.execute("SELECT COUNT(*) FROM books").fetchone()[0]
            if existing_count == 0:
                conn.executemany(
                    "INSERT INTO books (book_id, title, price_cents, available) VALUES (?, ?, ?, ?)",
                    [
                        ("b-100", "Distributed Systems for Beginners", 2999, 1),
                        ("b-101", "Go Concurrency in Practice", 3599, 1),
                        ("b-102", "Python Asyncio by Example", 3299, 0),
                    ],
                )
            conn.commit()

    def get_book(self, book_id: str) -> Book | None:
        with self._connect() as conn:
            row = conn.execute(
                "SELECT book_id, title, price_cents, available FROM books WHERE book_id = ?",
                (book_id,),
            ).fetchone()

        if row is None:
            return None

        return Book(
            book_id=row[0],
            title=row[1],
            price_cents=row[2],
            available=bool(row[3]),
        )

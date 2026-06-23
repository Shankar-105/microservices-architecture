from app.repository.catalog_repo import Book, CatalogRepository


class CatalogService:
    def __init__(self, repo: CatalogRepository) -> None:
        if repo is None:
            raise ValueError("catalog repository is required")
        self._repo = repo

    def health(self) -> str:
        return "ok"

    def get_book(self, book_id: str) -> Book | None:
        if not isinstance(book_id, str):
            return None
        normalized = book_id.strip()
        if not normalized or len(normalized) > 128:
            return None
        if any(ord(character) < 32 for character in normalized):
            return None
        return self._repo.get_book(normalized)

    def get_all_books(self) -> list[Book]:
        return self._repo.list_books()

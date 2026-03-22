from app.repository.catalog_repo import Book, CatalogRepository


class CatalogService:
    def __init__(self, repo: CatalogRepository) -> None:
        self._repo = repo

    def health(self) -> str:
        return "ok"

    def get_book(self, book_id: str) -> Book | None:
        return self._repo.get_book(book_id)

    def get_all_books(self) -> list[Book]:
        return self._repo.list_books()

import tempfile

import pytest

from app.repository.catalog_repo import CatalogRepository
from app.service.catalog_service import CatalogService


def test_catalog_service_normalizes_lookup():
    with tempfile.TemporaryDirectory() as directory:
        repo = CatalogRepository(f"{directory}/catalog.db")
        service = CatalogService(repo)
        book = service.get_book(" b-100 ")

    assert book is not None
    assert book.book_id == "b-100"


@pytest.mark.parametrize(
    "book_id",
    ["", "   ", "\x00b-100", "x" * 129],
)
def test_catalog_service_rejects_malformed_ids(book_id):
    with tempfile.TemporaryDirectory() as directory:
        repo = CatalogRepository(f"{directory}/catalog.db")
        service = CatalogService(repo)
        assert service.get_book(book_id) is None


def test_catalog_repository_rejects_empty_database_path():
    with pytest.raises(ValueError):
        CatalogRepository("")

import os

import pytest

from app.config import positive_port, required_db_path


@pytest.mark.parametrize("raw", ["0", "-1", "65536", "not-a-number"])
def test_positive_port_rejects_invalid_values(monkeypatch, raw):
    monkeypatch.setenv("TEST_PORT", raw)
    with pytest.raises(RuntimeError):
        positive_port("TEST_PORT", 50053)


def test_positive_port_accepts_boundary_values(monkeypatch):
    monkeypatch.setenv("TEST_PORT", "65535")
    assert positive_port("TEST_PORT", 50053) == 65535


@pytest.mark.parametrize("value", ["", "   ", "bad\x00path"])
def test_required_db_path_rejects_invalid_values(monkeypatch, value):
    monkeypatch.setenv("TEST_DB_PATH", value)
    with pytest.raises(RuntimeError):
        required_db_path("TEST_DB_PATH", "catalog.db")

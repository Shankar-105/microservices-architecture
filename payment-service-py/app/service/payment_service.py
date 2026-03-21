import uuid
from dataclasses import dataclass

from app.repository.payment_repo import PaymentRepository


@dataclass
class PaymentDecision:
    approved: bool
    transaction_id: str
    reason: str


class PaymentService:
    def __init__(self, repo: PaymentRepository) -> None:
        self._repo = repo

    def health(self) -> str:
        return "ok"

    def authorize(self, order_id: str, amount_cents: int) -> PaymentDecision:
        # Deterministic failure path for resilience testing: any order_id with prefix fail- is declined.
        if order_id.startswith("fail-"):
            decision = PaymentDecision(approved=False, transaction_id="", reason="deterministic test decline")
        elif amount_cents <= 0:
            decision = PaymentDecision(approved=False, transaction_id="", reason="invalid amount")
        else:
            txn_id = f"txn-{uuid.uuid4().hex[:12]}"
            decision = PaymentDecision(approved=True, transaction_id=txn_id, reason="approved")

        self._repo.create(order_id, amount_cents, decision.approved, decision.transaction_id, decision.reason)
        return decision

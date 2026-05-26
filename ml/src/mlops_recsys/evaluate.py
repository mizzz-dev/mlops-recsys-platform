from __future__ import annotations

from .data import Event


def precision_at_k(events: list[Event], recommendations: list[dict[str, object]], k: int = 10) -> float:
    if k <= 0:
        raise ValueError("k must be positive")
    relevant = {event.content_id for event in events if event.event_type in {"complete", "like"}}
    if not relevant:
        return 0.0
    predicted = {str(item["content_id"]) for item in recommendations[:k]}
    return len(predicted & relevant) / min(k, len(predicted) or k)

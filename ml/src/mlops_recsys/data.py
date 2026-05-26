from __future__ import annotations

import csv
import random
from collections import Counter, defaultdict
from dataclasses import dataclass
from pathlib import Path

EVENT_TYPES = ("view", "start", "complete", "like")


@dataclass(frozen=True)
class Event:
    user_id: str
    content_id: str
    event_type: str
    timestamp: str


def generate_events(path: Path, users: int = 100, events: int = 1000, seed: int = 42) -> None:
    random.seed(seed)
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", newline="", encoding="utf-8") as file:
        writer = csv.DictWriter(file, fieldnames=["user_id", "content_id", "event_type", "timestamp"])
        writer.writeheader()
        for index in range(events):
            user_number = random.randint(1, users)
            content_number = random.choices(range(1, 21), weights=list(range(20, 0, -1)), k=1)[0]
            writer.writerow({
                "user_id": f"user_{user_number:03d}",
                "content_id": f"quest_{content_number:03d}",
                "event_type": random.choice(EVENT_TYPES),
                "timestamp": f"2026-05-{(index % 28) + 1:02d}T12:00:00Z",
            })


def load_events(path: Path) -> list[Event]:
    with path.open("r", newline="", encoding="utf-8") as file:
        reader = csv.DictReader(file)
        events = [Event(**row) for row in reader]
    validate_events(events)
    return events


def validate_events(events: list[Event]) -> None:
    if not events:
        raise ValueError("events must not be empty")
    for event in events:
        if not event.user_id.startswith("user_"):
            raise ValueError(f"invalid user_id: {event.user_id}")
        if not event.content_id.startswith("quest_"):
            raise ValueError(f"invalid content_id: {event.content_id}")
        if event.event_type not in EVENT_TYPES:
            raise ValueError(f"invalid event_type: {event.event_type}")


def build_features(events: list[Event]) -> dict[str, object]:
    content_counts: Counter[str] = Counter()
    user_history: dict[str, set[str]] = defaultdict(set)
    for event in events:
        weight = 2 if event.event_type in {"complete", "like"} else 1
        content_counts[event.content_id] += weight
        user_history[event.user_id].add(event.content_id)

    total = sum(content_counts.values()) or 1
    popularity = [
        {"content_id": content_id, "score": count / total, "reason": "popular_content"}
        for content_id, count in content_counts.most_common()
    ]
    return {"content_popularity": popularity, "user_count": len(user_history), "event_count": len(events)}

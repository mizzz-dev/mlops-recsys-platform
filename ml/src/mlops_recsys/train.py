from __future__ import annotations

import argparse
import json
from datetime import UTC, datetime
from pathlib import Path

from .data import build_features, generate_events, load_events
from .evaluate import precision_at_k


def train(events_path: Path, model_path: Path) -> dict[str, object]:
    if not events_path.exists():
        generate_events(events_path)

    events = load_events(events_path)
    features = build_features(events)
    recommendations = features["content_popularity"][:20]
    precision = precision_at_k(events, recommendations, k=10)

    model = {
        "version": "popular-baseline-local",
        "generated_at": datetime.now(UTC).isoformat(),
        "model_type": "popular_baseline",
        "evaluation": {
            "precision_at_10": round(precision, 4),
            "event_count": features["event_count"],
            "user_count": features["user_count"],
        },
        "recommendations": recommendations,
    }

    model_path.parent.mkdir(parents=True, exist_ok=True)
    model_path.write_text(json.dumps(model, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return model


def main() -> None:
    parser = argparse.ArgumentParser(description="Train local recommendation baseline model.")
    parser.add_argument("--events", type=Path, default=Path("data/samples/events.csv"))
    parser.add_argument("--model", type=Path, default=Path("artifacts/model.json"))
    args = parser.parse_args()

    model = train(args.events, args.model)
    print(json.dumps({"model_path": str(args.model), "evaluation": model["evaluation"]}, ensure_ascii=False))


if __name__ == "__main__":
    main()

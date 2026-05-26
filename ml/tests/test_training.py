from __future__ import annotations

import json

from mlops_recsys.data import build_features, generate_events, load_events
from mlops_recsys.evaluate import precision_at_k
from mlops_recsys.train import train


def test_generate_and_load_events(tmp_path):
    events_path = tmp_path / "events.csv"
    generate_events(events_path, users=3, events=10, seed=1)

    events = load_events(events_path)

    assert len(events) == 10
    assert events[0].user_id.startswith("user_")
    assert events[0].content_id.startswith("quest_")


def test_build_features_returns_popularity(tmp_path):
    events_path = tmp_path / "events.csv"
    generate_events(events_path, users=5, events=50, seed=2)
    events = load_events(events_path)

    features = build_features(events)

    assert features["event_count"] == 50
    assert features["user_count"] > 0
    assert len(features["content_popularity"]) > 0


def test_precision_at_k_is_bounded(tmp_path):
    events_path = tmp_path / "events.csv"
    generate_events(events_path, users=5, events=50, seed=3)
    events = load_events(events_path)
    features = build_features(events)

    score = precision_at_k(events, features["content_popularity"], k=5)

    assert 0.0 <= score <= 1.0


def test_train_writes_model_artifact(tmp_path):
    events_path = tmp_path / "events.csv"
    model_path = tmp_path / "model.json"

    model = train(events_path, model_path)

    assert model_path.exists()
    saved = json.loads(model_path.read_text(encoding="utf-8"))
    assert saved["version"] == model["version"]
    assert len(saved["recommendations"]) > 0
    assert "precision_at_10" in saved["evaluation"]

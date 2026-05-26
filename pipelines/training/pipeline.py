from __future__ import annotations

import argparse
import json
from pathlib import Path

PIPELINE_SPEC = {
    "name": "local-recsys-training-pipeline",
    "description": "Local compile target for the recommendation training workflow.",
    "components": [
        {"name": "ingest", "outputs": ["events"]},
        {"name": "validate", "inputs": ["events"], "outputs": ["validated_events"]},
        {"name": "feature_build", "inputs": ["validated_events"], "outputs": ["features"]},
        {"name": "train", "inputs": ["features"], "outputs": ["model"]},
        {"name": "evaluate", "inputs": ["model", "validated_events"], "outputs": ["metrics"]},
        {"name": "register", "inputs": ["model", "metrics"], "outputs": ["registered_model"]},
        {"name": "batch_predict", "inputs": ["registered_model"], "outputs": ["recommendations"]}
    ]
}


def compile_pipeline(output: Path) -> None:
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(PIPELINE_SPEC, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser(description="Compile local pipeline spec.")
    parser.add_argument("--compile", action="store_true", help="Compile pipeline spec")
    parser.add_argument("--output", type=Path, default=Path("artifacts/pipeline.json"))
    args = parser.parse_args()

    if not args.compile:
        raise SystemExit("--compile is required")
    compile_pipeline(args.output)
    print(json.dumps({"pipeline_path": str(args.output)}, ensure_ascii=False))


if __name__ == "__main__":
    main()

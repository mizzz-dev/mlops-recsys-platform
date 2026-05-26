# 0001: Use Go for inference API

## Decision

Use Go for the inference API.

## Context

The serving layer should be small, easy to deploy as a container, and simple to operate. The training layer remains in Python, so using Go for serving makes the boundary between training and serving explicit.

## Consequences

- The API can be tested with the Go standard library.
- The runtime dependency surface stays small.
- Model artifacts must use a language-neutral format such as JSON for the MVP.

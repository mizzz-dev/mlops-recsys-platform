---
name: Operations task
description: Track an operational improvement or deployment follow-up
title: "[Ops] "
labels: ["ops"]
body:
  - type: textarea
    id: purpose
    attributes:
      label: Purpose
      description: What operational problem does this task solve?
    validations:
      required: true
  - type: textarea
    id: scope
    attributes:
      label: Scope
      description: What is included and excluded?
    validations:
      required: true
  - type: textarea
    id: verification
    attributes:
      label: Verification
      description: How will this be verified safely?
    validations:
      required: true
  - type: textarea
    id: rollback
    attributes:
      label: Rollback plan
      description: How can this be rolled back?
    validations:
      required: false
---

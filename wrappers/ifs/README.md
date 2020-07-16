# If wrapper

## Purpose
Conditionally execute a node based on conditions.
If multiple conditions are passed, all conditions must be validated.

## YAML
```yaml
nodes:
  - id: node_id
    if:
      - conditionA
      - conditionB
```

## Implemented interface methods
Nodes:
- Inputs
- Run

# After

This wrapper ensures that a node gets executed only once all the nodes listed in its `after` key word have been executed.

```yaml
id: AfterExample

nodes:
  - id: node_a

  - id: node_b

  # node_c will start only once node_a and node_b have finished
  - id: node_c
    after:
      - node_a
      - node_b
```

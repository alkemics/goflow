# goflow

goflow is a workflow tool based on code generation. Its goal is to make it easy to write and maintain workflows.

The workflows are represented as DAGs. The nodes of a DAG are written in Go and the edges are written in a YAML file. Goflow will then generate the orchestration code based on the Go nodes and the YAML edges.

Head over to the [example](https://github.com/alkemics/goflow/tree/master/example) directory to see how to use it.

## Design

Because goflow is made to make it easier to write and maintain workflows, its design tries to also make it easy to maintain and extend the tool itself. The basic loader knows only very few keywords: `name`, `package` and `nodes`. Every other feature is then added via a `GraphWrapper`.

## Caveats / future improvements

- the wrappers are not all independent, meaning that you need to be careful about the order of the wrappers passed to the `Load` function. This could be solved by introducing a resolver for those dependencies.
- it is still Go oriented: the wrappers generating code already present in the `wrappers` directory generate Go code. Moving the code generating function of those wrappers to the `gfutil/gfgo` package would make sense given that one of goflow's ambitions is to become language agnostic. Howevere, more exploration is needed first to better decide how to do that, for example by trying goflow with Python or JS.

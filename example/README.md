# goflow example

To generate the graphs, run:

```bash
go run ./generator/main.go
```

You can then start the playground by running:

```bash
go build -o goflowpg ./main.go
./goflowpg
```

You can then try the playground with:

```bash
$ curl 'http://127.0.0.1:8080'
# ... list of the graphs

$ curl 'http://127.0.0.1:8080?name=github.com/alkemics/goflow/example/graphs.Affiner' -d '{"a": 3, "x": 2, "b": 5}'
{"result":11}
```

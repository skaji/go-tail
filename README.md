# go tail

WIP

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
ch := tail.Run(ctx, tail.NewConfig("access.log"))
for line := range ch {
    fmt.Printf("got [%s]\n", line)
}
```

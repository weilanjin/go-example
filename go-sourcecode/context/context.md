# context

Package context src/context/context.go

emptyCtx 是 parent context

WithTimeout 用法

```go
	ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
	defer cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	}

	// Output:
	// context deadline exceeded
```
# go-matcha

go-matcha is a simple matcher for Go. It provides the flexible matching of various value and type.

## How to Use

```go

matcha.Equal(1, 1) // true
matcha.Equal(matcher.BeAny(), 1)           // any value is expected and return "true"
matcha.Equal(matcher.BeInt(), 1)           // int is expected and return "true"
matcha.Equal(matcher.BeString(), 1)        // string is expected and return "false"
matcha.Equal(matcher.BeInt().Not(), 1)     // not int is expected and return "false"
matcha.Equal(matcher.BeInt().Pointer(), 1) // pointer int is expected and return "false"
```

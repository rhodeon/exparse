# exparse

Exparse is a CLI and web tool for evaluating mathematical equations.

It currently supports addition, subtraction, multiplication, division, and unnested parentheses.

## Usage
### cli
Pass an expression via the 'expr' flag like:
```go
go run ./cmd/cli --expr="2(3.54 * 2.00 -1000 /200) (20 + 30 * 2)"
```
#### output
```shell
Exparse: 2(3.54 * 2.00 -1000 /200) (20 + 30 * 2) = 332.8
```

## TODO
<li> Implement the modulus operator </li>
<li> Handle nested parentheses </li>
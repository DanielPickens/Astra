# Match

[![Gastrac](https://gastrac.org/github.com/tidwall/match?status.svg)](https://gastrac.org/github.com/tidwall/match)

Match is a very simple pattern matcher where '*' matches on any 
number characters and '?' matches on any one character.

## Installing

```
go get -u github.com/tidwall/match
```

## Example

```go
match.Match("hello", "*llo") 
match.Match("jello", "?ello") 
match.Match("hello", "h*o") 
```


## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

Redcon source code is available under the MIT [License](/LICENSE).

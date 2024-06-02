# json5b 

Now based on `titanous/json5` package instead of `yosuke-furukawa/json5` [![GoDoc](https://godoc.org/github.com/titanous/json5?status.svg)](https://godoc.org/github.com/titanous/json5) [![Build Status](https://github.com/titanous/json5/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/titanous/json5/actions/workflows/ci.yml)

This is a Go package that implements decoding of
[JSON5](https://github.com/json5/json5). See [the
documentation](https://godoc.org/github.com/titanous/json5) for usage information.

- The tag being used is `json5` instead of `json`
- merged the patch from `skybosi` to support autoconvert string to number, added more tests



# HOW TO USE

```
go install github.com/kokizzu/json5b 
$ json5 -c path/to/test.json5 # output stdout
$ json5 -c path/to/test.json5 -o path/to/test.json # output path/to/test.json
```

# go get
```
$ go get github.com/kokizzu/json5b
```

# example

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kokizzu/json5b/encoding/json5b"
)

func main() {
	var data interface{}
	dec := json5b.NewDecoder(os.Stdin)
	err := dec.Decode(&data)
	if err != nil {
		fmt.Println(err)
	}
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}
```

```js
// This is json5 demo
// json5 can write comment in your json
{
  key : "Key does not need double quote",
  // json specific
  "of" : "course we can use json as json5",
  trailing : "trailing comma is ok",
}
```

```
$ json5 -c example.json5
# output
#{
#    "key": "Key does not need double quote",
#    "of": "course we can use json as json5",
#    "trailing": "trailing comma is ok"
#}
```

# TODO
- block comment
- multiline string
- hexadecimal notation




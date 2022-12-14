json5 implemented by golang
================

Forked version of [yosuke-furukawa/json5](https://github.com/yosuke-furukawa/json5), uses `json5` tag instead of default `json` tag.

[JSON5](https://github.com/aseemk/json5) is Yet another JSON.

# INSTALL

```
$ brew tap yosuke-furukawa/json5 # original one
$ brew install json5
```

# HOW TO USE

```
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
	"github.com/kokizzu/json5b/encoding/json5b"
	"os"
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




# Goman

## Description
Retrieves documentation of go standard from https://pkg.go.dev/ for a specified
function/constant/type and displays its name, signature, description and usage
example if present.

## Installation
```
$ go build .
```
## Usage
```
$ goman fmt.Println
Package fmt
Documentation: https://pkg.go.dev/fmt@go1.23.1
Package fmt implements formatted I/O with functions analogous to C's printf and scanf.
---
function Println
func Println(a ...any) (n int, err error)
---
Println formats using the default formats for its operands and writes to standard output.
Spaces are always added between operands and a newline is appended.
It returns the number of bytes written and any write error encountered.

package main

import (
    "fmt"
)

func main() {
    const name, age = "Kim", 22
    fmt.Println(name, "is", age, "years old.")

    // It is conventional not to worry about any
    // error returned by Println.

}
```
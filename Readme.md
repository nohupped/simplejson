# { SimpleJSON } [![Build Status](https://travis-ci.com/nohupped/simplejson.svg?branch=master)](https://travis-ci.com/nohupped/simplejson) [![codecov.io](https://codecov.io/github/nohupped/simplejson/coverage.svg?branch=master)](https://codecov.io/github/nohupped/simplejson?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/nohupped/simplejson)](https://goreportcard.com/report/github.com/nohupped/simplejson)

A lame attempt to re-create python's `json.loads()` using `Get()` methods to retrive a key by moving the ugly typeassertions to the library.

## Why

This is just a wrapper over go's `encoding/json`. This uses an `interface{}` to unmarshal the json and the exposed interface method `Get()` uses `type assertions` internally to get the value of a key if it is a hashmap or the index if it is an array. Since this moves the ugly typeassertions for nested json parsing inside the library, the code becomes easier to read.

## Example

Consider the json string

```go
sampleJSON := `{
    "Actors":
    [
        {
            "name": "Tom Cruise",
            "age": 56,
            "Born At": "Syracuse, NY",
            "Birthdate": "July 3, 1962",
            "hasChildren":true,
            "hasTwitterAccount":true,
            "hasGreyHair":false,
            "photo": "https://jsonformatter.org/img/tom-cruise.jpg"
        },
        {
            "name": "Robert Downey Jr.",
            "age": 53,
            "Born At": "New York City, NY",
            "Birthdate": "April 4, 1965",
            "photo": "https://jsonformatter.org/img/Robert-Downey-Jr.jpg"
        }
    ]
}`
```

This can be parsed and values can be fetched as

```go

d, err := simplejson.Loads(sampleJSON)
if err != nil {
    panic(err)
}
data := d.Get("Actors").Get("", 0).Get("name")
fmt.Println(data) // outputs Tom Cruise

```

Consider the json file `samplejson.json` in this repository.

This file can be parsed and the values can be fetched as

```go
fd, err := os.Open("samplejson.json")
if err != nil {
    panic(err)
}
d1, err := Load(fd)
if err != nil{
    panic(err)
}
t.Logf("%s", d1.Get("", 2).Get("tags").Get("", 3))

fd.Close()

```

It is the responsibility of the caller to close the file descriptor.

Because this uses `interface{}` to unmarshal json, the structure of json is not required to be predefined. Each `Get()` evaluates the type and extracts accordingly.

## `Get()` function ~~panics~~ returns `nil` for the method `Bytes()` if the key is not present

Instead of returning an error, `Get()` function ~~panics~~ returns `nil` for the `Bytes()` method when trying to fetch an invalid key or trying to read from an invalid index. This is done to easily chain the method as follows

```go
Get("foo").Get("bar").Get("", 2)
```

without having to evaluate for errors every time or to have to write a `recover` as below.

```go
defer func() {
    if r := recover(); r != nil {...}
```

Evaluate for `nil` value before using it.

Eg:

```go
d, err := simplejson.Loads([]byte(sampleJSON))
if err != nil {
    ...
}
d1 := d.Get("Actors").Get("", 0).Get("names").Get("foo").Get("Bar")
if d1.Bytes != nil {
    ...
}
...

```

## `String()` function returns empty string `""` for `nil` values

This is done to avoid the

```bash
panic: runtime error: invalid memory address or nil pointer dereference
```

at runtime. Evaluate for `nil` value before applying `String()` method.

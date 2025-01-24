+++
title = "hello there"
description = "the quick and easy markdown file to see progress"
date = 2025-01-24

[author]
name = "Devin Ward"
email = "devin@hey.com"
+++

this is a check for markdown. a simple test file

## list
 - step 1
 - step 2

## code

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

::: Aside Title
This will be an aside.

Asides can include all normal markdown too like lists and code.

- [x] Step 1
- [ ] Step 2

```go
type Doer interface {
    Do() error
}
```

:::

Additional content outside aside.
# Scipio

Static website generator in Golang

![Go](https://github.com/lchsk/scipio/workflows/Go/badge.svg)

[Blog article](https://lchsk.com/static-websites-in-golang-and-rust.html)

---

## List of available template variables

```
{{description}}

{{keywords}}

{{title}}

{{posts-begin}}
{{post_link}}
{{post_date}}
{{posts-end}}

{{@<page-slug>}}

{{#include <template.html>}}


```

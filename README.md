# Scipio

Static website generator in Golang

![Go](https://github.com/lchsk/scipio/workflows/Go/badge.svg)

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

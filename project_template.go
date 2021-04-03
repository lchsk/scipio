package main

const firstPostPattern = `
---
title: Welcome to your new website
created: 2021-01-27T00:00:00Z
description: Description of the first post
keywords: keyword1, keyword2, etc.
---

### Welcome

__Source code example__
`

const privacyPolicy = `
---
title: Privacy policy
created: 2021-01-27T00:00:00Z
description: Privacy policy
keywords: privacy policy
---

## Privacy policy
`

const indexPage = `
---
title: Title
created: 2021-01-27T00:00:00Z
description: Description
keywords: keyword1, etc.
---
`

const headerTheme = `
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" type="text/css" href="/static/styles.css">
    <link rel="alternate" type="application/rss+xml" href="/posts.xml" />
`
const topTheme = `
	<a href="/">Home</a>
`
const footerTheme = `
	<div>All rights reserved</div>
`

const indexTheme = `
<!DOCTYPE html>
  <head>
    <meta charset="UTF-8">
    <meta name="description" content="{{description}}">
    <meta name="keywords" content="{{keywords}}">
    <meta name="author" content="">
    <title>{{title}} - {{description}}</title>
    {{#include header.html}}
  </head>
  <body>
  <div class="container">
  {{#include top.html}}

  <main role="main" class="container">
    <div class="row">
      <div class="col-md-12 blog-main">
        {{posts-begin}}
        <div class="blog-post">
        <h2>{{post_link}}</h2>
	      <p>{{post_date}}</p>
          <p>{{post_description}}</p>
        </div>
        {{posts-end}}
      </div>
    </div>
  </main>

  {{#include footer.html}}
  </body>
</html>
`

const postTheme = `
<!DOCTYPE html>
  <head>
    <meta charset="UTF-8">
    <meta name="description" content="{{description}}">
    <meta name="keywords" content="{{keywords}}">
    <meta name="author" content="">
    <title>{{title}}</title>
    {{#include header.html}}
  </head>
  <body>
  <div class="container">
  {{#include top.html}}

  <main role="main" class="container">
    <div class="row">
      <div class="col-md-12 blog-main">
        <div class="blog-post">
          <h2>{{title}}</h2>
          <p>{{date}}</p>
          {{body}}
        </div>
      </div>
    </div>
  </main>

  {{#include footer.html}}
  </div>
  </body>
</html>
`

const pageTheme = `
<!DOCTYPE html>
  <head>
    <meta charset="UTF-8">
    <meta name="description" content="{{description}}">
    <meta name="keywords" content="{{keywords}}">
    <meta name="author" content="">
    <title>{{title}}</title>
    {{#include header.html}}
  </head>
  <body>
  <div class="container">
  {{#include top.html}}

    <main role="main" class="container">
      <div class="row">
        <div class="col-md-12 blog-main">
          <div class="blog-post">
            <h2>{{title}}</h2>
            {{body}}
          </div>
        </div>
      </div>
    </main>

    {{#include footer.html}}
  </div>
  </body>
</html>
`

const postsTheme = `
<!DOCTYPE html>
  <head>
    <meta charset="UTF-8">
    <meta name="description" content="{{description}}">
    <meta name="keywords" content="{{keywords}}">
    <meta name="author" content="">
    <title>{{title}}</title>
    {{#include header.html}}
  </head>
  <body>
  <div class="container">
    {{#include top.html}}

    <main role="main" class="container">
      <div class="row">
        <div class="col-md-12 blog-main">
          <div class="blog-post">
            <ul id="post-list">
              {{posts-begin}}
              <li>{{post_link}}</li>
              <p>{{post_description}}</p>
              {{posts-end}}
            </ul>
          </div>
        </div>
      </div>
    </main>

    {{#include footer.html}}
  </div>
  </body>
</html>
`

const themeStyleApp = `
.blog-header {
  line-height: 1;
}

.blog-header-logo {
  font-family: sans-serif;
  font-size: 2.25rem;
}

.blog-header-logo:hover {
  text-decoration: none;
}

h1, h2, h3, h4, h5, h6 {
  font-family: sans-serif;
}

.display-4 {
  font-size: 2.5rem;
}
@media (min-width: 768px) {
  .display-4 {
    font-size: 3rem;
  }
}

.nav-scroller {
  position: relative;
  z-index: 2;
  height: 2.75rem;
  overflow-y: hidden;
}

.nav-scroller .nav {
  display: -ms-flexbox;
  display: flex;
  -ms-flex-wrap: nowrap;
  flex-wrap: nowrap;
  padding-bottom: 1rem;
  margin-top: -1px;
  overflow-x: auto;
  text-align: center;
  white-space: nowrap;
  -webkit-overflow-scrolling: touch;
}

.nav-scroller .nav-link {
  padding-top: .75rem;
  padding-bottom: .75rem;
  font-size: .875rem;
}

.card-img-right {
  height: 100%;
  border-radius: 0 3px 3px 0;
}

.flex-auto {
  -ms-flex: 0 0 auto;
  flex: 0 0 auto;
}

.h-250 { height: 250px; }
@media (min-width: 768px) {
  .h-md-250 { height: 250px; }
}

/* Pagination */
.blog-pagination {
  margin-bottom: 4rem;
}
.blog-pagination > .btn {
  border-radius: 2rem;
}

/*
 * Blog posts
 */
.blog-post {
  margin-bottom: 4rem;
}
.blog-post-title {
  margin-bottom: .25rem;
  font-size: 2.5rem;
}
.blog-post-meta {
  margin-bottom: 1.25rem;
  color: #999;
}

/*
 * Footer
 */
.blog-footer {
  padding: 2.5rem 0;
  color: #555;
  text-align: left;
}
.blog-footer p:last-child {
  margin-bottom: 0;
}

img.rounded {
  @extend .img-fluid
}

img {
  @extend .rounded;
}

code {
  color: $primary;
}

/*
 * Pygments style for source code syntax highlighting
 */

.highlight .hll { background-color: #ffffcc }
.highlight  { background: #f8f8f8; }
.highlight .c { color: #008800; font-style: italic } /* Comment */
.highlight .err { border: 1px solid #FF0000 } /* Error */
.highlight .k { color: #AA22FF; font-weight: bold } /* Keyword */
.highlight .o { color: #666666 } /* Operator */
.highlight .ch { color: #008800; font-style: italic } /* Comment.Hashbang */
.highlight .cm { color: #008800; font-style: italic } /* Comment.Multiline */
.highlight .cp { color: #008800 } /* Comment.Preproc */
.highlight .cpf { color: #008800; font-style: italic } /* Comment.PreprocFile */
.highlight .c1 { color: #008800; font-style: italic } /* Comment.Single */
.highlight .cs { color: #008800; font-weight: bold } /* Comment.Special */
.highlight .gd { color: #A00000 } /* Generic.Deleted */
.highlight .ge { font-style: italic } /* Generic.Emph */
.highlight .gr { color: #FF0000 } /* Generic.Error */
.highlight .gh { color: #000080; font-weight: bold } /* Generic.Heading */
.highlight .gi { color: #00A000 } /* Generic.Inserted */
.highlight .go { color: #888888 } /* Generic.Output */
.highlight .gp { color: #000080; font-weight: bold } /* Generic.Prompt */
.highlight .gs { font-weight: bold } /* Generic.Strong */
.highlight .gu { color: #800080; font-weight: bold } /* Generic.Subheading */
.highlight .gt { color: #0044DD } /* Generic.Traceback */
.highlight .kc { color: #AA22FF; font-weight: bold } /* Keyword.Constant */
.highlight .kd { color: #AA22FF; font-weight: bold } /* Keyword.Declaration */
.highlight .kn { color: #AA22FF; font-weight: bold } /* Keyword.Namespace */
.highlight .kp { color: #AA22FF } /* Keyword.Pseudo */
.highlight .kr { color: #AA22FF; font-weight: bold } /* Keyword.Reserved */
.highlight .kt { color: #00BB00; font-weight: bold } /* Keyword.Type */
.highlight .m { color: #666666 } /* Literal.Number */
.highlight .s { color: #BB4444 } /* Literal.String */
.highlight .na { color: #BB4444 } /* Name.Attribute */
.highlight .nb { color: #AA22FF } /* Name.Builtin */
.highlight .nc { color: #0000FF } /* Name.Class */
.highlight .no { color: #880000 } /* Name.Constant */
.highlight .nd { color: #AA22FF } /* Name.Decorator */
.highlight .ni { color: #999999; font-weight: bold } /* Name.Entity */
.highlight .ne { color: #D2413A; font-weight: bold } /* Name.Exception */
.highlight .nf { color: #00A000 } /* Name.Function */
.highlight .nl { color: #A0A000 } /* Name.Label */
.highlight .nn { color: #0000FF; font-weight: bold } /* Name.Namespace */
.highlight .nt { color: #008000; font-weight: bold } /* Name.Tag */
.highlight .nv { color: #B8860B } /* Name.Variable */
.highlight .ow { color: #AA22FF; font-weight: bold } /* Operator.Word */
.highlight .w { color: #bbbbbb } /* Text.Whitespace */
.highlight .mb { color: #666666 } /* Literal.Number.Bin */
.highlight .mf { color: #666666 } /* Literal.Number.Float */
.highlight .mh { color: #666666 } /* Literal.Number.Hex */
.highlight .mi { color: #666666 } /* Literal.Number.Integer */
.highlight .mo { color: #666666 } /* Literal.Number.Oct */
.highlight .sa { color: #BB4444 } /* Literal.String.Affix */
.highlight .sb { color: #BB4444 } /* Literal.String.Backtick */
.highlight .sc { color: #BB4444 } /* Literal.String.Char */
.highlight .dl { color: #BB4444 } /* Literal.String.Delimiter */
.highlight .sd { color: #BB4444; font-style: italic } /* Literal.String.Doc */
.highlight .s2 { color: #BB4444 } /* Literal.String.Double */
.highlight .se { color: #BB6622; font-weight: bold } /* Literal.String.Escape */
.highlight .sh { color: #BB4444 } /* Literal.String.Heredoc */
.highlight .si { color: #BB6688; font-weight: bold } /* Literal.String.Interpol */
.highlight .sx { color: #008000 } /* Literal.String.Other */
.highlight .sr { color: #BB6688 } /* Literal.String.Regex */
.highlight .s1 { color: #BB4444 } /* Literal.String.Single */
.highlight .ss { color: #B8860B } /* Literal.String.Symbol */
.highlight .bp { color: #AA22FF } /* Name.Builtin.Pseudo */
.highlight .fm { color: #00A000 } /* Name.Function.Magic */
.highlight .vc { color: #B8860B } /* Name.Variable.Class */
.highlight .vg { color: #B8860B } /* Name.Variable.Global */
.highlight .vi { color: #B8860B } /* Name.Variable.Instance */
.highlight .vm { color: #B8860B } /* Name.Variable.Magic */
.highlight .il { color: #666666 } /* Literal.Number.Integer.Long */
`

const themeStyleBootstrap = `
$my-color-1: #0000ff;
$my-color-2: #ff00ff;
$my-color-3: #000000;
$my-color-4: #555555;
$my-color-5: #ffff00;

$primary: $my-color-1;
$secondary: $my-color-2;
$light: $my-color-3;
$dark: $my-color-4;

$btn-padding-y: 5px;
$btn-padding-x: 5px;
$btn-border-radius: 5px;
$btn-border-width: 0;
$input-btn-focus-width: 0;

@import "../../node_modules/bootstrap/scss/bootstrap";
`

const themeStyleBundle = `
@import "bootstrap";
@import "app";
`

const configValue = `
# Base URL of the website
url = "https://www.mywebsite.com"

output_extension = ".html"
links_begin_with_slash = false

# Settings for the RSS feed
[rss]
generate_rss = true
title = "My Website"
description = "Interesting new website"
author_name = ""
author_email = ""

# Settings for the posts.html listing all the posts
[posts]
slug = "posts"
title = "Posts"
Description = "Interesting new website"
Keywords = ["website", "hello"]

[static]
copy = [
    # {from = "sitemap.xml", to = "./"},
]
`

const packageJsonValue = `
{
  "name": "mywebsite",
  "version": "1.0.0",
  "description": "",
  "main": "",
  "scripts": {
    "sass": "./node_modules/sass/sass.js --no-source-map --style=compressed ./themes/default/bundle.scss ./themes/default/static/styles.css",
    "sass-watch": "./node_modules/sass/sass.js --no-source-map --style=compressed --watch ./themes/default/bundle.scss ./themes/default/static/styles.css"
  },
  "keywords": [],
  "author": "",
  "license": "",
  "devDependencies": {
    "bootstrap": "^5.0.0-beta3",
    "sass": "^1.32.8"
  }
}
`
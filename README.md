# Migrate Geek Igor

`migrate_geek_igor` is a helper tool to _migrate_ my old jekyll based blog + content to new structure.

It is intended as a single use tool.

Usage:
```src
migrate_geek_igor -from ~/jekyll_blog  -to ./content
```

## Content structure

Old content structure:
```
.
├── _posts
│   ├── 2013-03-22-meta-blogging-jekyll-setup.md
│   └── 2022-01-30-private-link.md
└── static
    └── img
    └── posts
        ├── 2022-01-ip-filtering.png
        ├── 2022-01-private-link-basic.png
        ├── 2022-01-private-link.png
        └── 2022-01-vpc-peering.png
```

New content structure:
```
.
├── _posts
│   └── 2022-01-30-private-link.md
├── _assets
│   ├── 2022-01-ip-filtering.png
│   ├── 2022-01-ip-filtering-basic.png
│   ├── 2022-01-private-link.png
│   └── 2022-01-vpc-peering.png
(...)
```


## Front matter

Old content front matter:
```
---
layout: post
title: "Private Link is the IP filtering of the cloud"
tags: [privatelink, network]
tagline: Use cases for Private Link and differences in its implementation across the major Cloud Providers
---
```

New content front matter:
```
---
tags:
- privatelink
- network
redirect_from:
- 2022/01/30/private-link.html
---
# Private Link is the IP filtering of the cloud

Use cases for Private Link and differences in its implementation across the major Cloud Providers.
```


## Markdown content

- Rewrite image URLs

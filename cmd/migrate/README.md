# Migrate

`migrate` is a helper tool to _migrate_ my old jekyll based blog + content to new structure.

It is intended as a single use tool.

Usage:
```src
migrate -from ~/jekyll_blog  -to ./content
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
├── 2022
│   ├── private-link.md
│   ├── ip-filtering.png
│   ├── ip-filtering-basic.png
│   ├── private-link.png
│   └── vpc-peering.png
└── archive
└── 2013
└── meta-blogging-jekyll-setup.md
```

Transformations to file system structure:
- [X] Put all content under a year (drop the _post / static distinction)
- [X] Drop date from file name
- [ ] Put some content under /archive, depending on a front matter setting


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
tags: [privatelink, network]
date: 2022-01-30
from:
    - 2022/01/30/private-link.html
---
# Private Link is the IP filtering of the cloud

Use cases for Private Link and differences in its implementation across the major Cloud Providers.
```

Transformations to front matter
- [ ] Drop `layout`
- [ ] Drop `title` (but add it to content)
- [ ] Drop `tagline` (but add it to content)
- [ ] Add `date` based on the prev file name
- [ ] Add `from` based on the url scheme of my blog


## Markdown content

Transformations to file content:
- [ ] Add title in the first line
- [ ] Add tagline after the first line
- [ ] Rewrite image urls


## Validation
- [ ] Eyeball top 10 pages
- [ ] Compare jekyll output to `from` property in the front matter

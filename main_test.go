package main

import (
	"reflect"
	"testing"
)

type fnameExample struct {
	fname      string
	title      string
	datePrefix string
	ext        string
}

var examples = []fnameExample{
	{fname: "2018-12-08-fire-and-motion.md", title: "fire-and-motion", datePrefix: "2018-12-08", ext: "md"},
	{fname: "2014-07-31-home-office.jpg", title: "home-office", datePrefix: "2014-07-31", ext: "jpg"},
	{fname: "2020-11-debugger.png", title: "debugger", datePrefix: "2020-11", ext: "png"},
	{fname: "2020-11-0-instances-runnning.png", title: "0-instances-runnning", datePrefix: "2020-11", ext: "png"},
}

func Test_splitFname(t *testing.T) {
	for _, tt := range examples {
		tt := tt
		t.Run(tt.fname, func(t *testing.T) {
			t.Parallel()
			gotDatePrefix, gotTitle, gotExt := splitFname(tt.fname)
			if gotDatePrefix != tt.datePrefix {
				t.Errorf("splitFname().datePrefix = %v, want %v", gotDatePrefix, tt.datePrefix)
			}
			if gotTitle != tt.title {
				t.Errorf("splitFname().title = %v, want %v", gotTitle, tt.title)
			}
			if gotExt != tt.ext {
				t.Errorf("splitFname().ext = %v, want %v", gotExt, tt.ext)
			}
		})
	}
}

func Test_processPost(t *testing.T) {
	type args struct {
		fname   string
		date    string
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Post
		wantErr bool
	}{
		{
			name: "Should process front matter and content",
			args: args{
				fname: "private-link",
				date:  "2022-01-30",
				content: []byte(`---
layout: post
title: "Private Link is the IP filtering of the cloud"
tags: [privatelink, network]
tagline: Use cases for Private Link and differences in its implementation across the major Cloud Providers.
---

foo
`),
			},
			want: &Post{
				Content: []byte(`---
layout: post
tags:
- privatelink
- network
date: "2022-01-30"
redirect_from:
- 2022/01/30/private-link.html
---
# Private Link is the IP filtering of the cloud

Use cases for Private Link and differences in its implementation across the major Cloud Providers.

foo
`),
			},
		},
		{
			name: "Should rewrite image locations",
			args: args{
				fname: "private-link",
				date:  "2022-01-30",
				content: []byte(`---
layout: post
title: "Private Link is the IP filtering of the cloud"
---
<center style="float: right; display: block; margin: 10px;">
	<img alt='Basic Personal Kanban' src='/static/img/posts/2013-08-29-pk-book.png' />
	<br/>
	<em>Img 1.</em> Personal Kanban Book
</center>

![Gospodarka, Głupcze](/static/img/posts/2018-04-gospodarka.jpg)

![Spacemacs](/static/img/posts/2018-04-spacemacs.png)

![Private Link](/static/img/posts/2022-01-private-link-basic.png)
`),
			},
			want: &Post{
				Content: []byte(`---
layout: post
tags: []
date: "2022-01-30"
redirect_from:
- 2022/01/30/private-link.html
---
# Private Link is the IP filtering of the cloud

<center style="float: right; display: block; margin: 10px;">
	<img alt='Basic Personal Kanban' src='pk-book.png' />
	<br/>
	<em>Img 1.</em> Personal Kanban Book
</center>

![Gospodarka, Głupcze](gospodarka.jpg)

![Spacemacs](spacemacs.png)

![Private Link](private-link-basic.png)
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processPost(tt.args.fname, tt.args.date, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("processPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Content, tt.want.Content) {
				t.Errorf("processPost().Content got = %v, want %v", string(got.Content), string(tt.want.Content))
			}
		})
	}
}

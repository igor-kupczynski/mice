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
		created string
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
				fname:   "private-link",
				created: "2022-01-30",
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
tags:
- privatelink
- network
created: "2022-01-30"
from:
- 2022/01/30/private-link.html
---
# Private Link is the IP filtering of the cloud

Use cases for Private Link and differences in its implementation across the major Cloud Providers.

foo
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processPost(tt.args.fname, tt.args.created, tt.args.content)
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

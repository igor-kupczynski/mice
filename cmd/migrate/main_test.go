package main

import "testing"

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

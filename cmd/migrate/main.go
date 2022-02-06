// migrate is a helper tool to _migrate_ my old jekyll based blog + content to new structure.
//
// It is intended as a single use tool.
//
// Usage:
//   migrate -from ~/jekyll_blog  -to ./content
//
//
// Old content structure:
// .
// ├── _posts
// │   ├── 2013-03-22-meta-blogging-jekyll-setup.md
// │   └── 2022-01-30-private-link.md
// └── static
//     └── img
//         └── posts
//             ├── 2022-01-ip-filtering.png
//             ├── 2022-01-private-link-basic.png
//             ├── 2022-01-private-link.png
//             └── 2022-01-vpc-peering.png
//
// New content structure:
// .
// ├── 2022
// │   ├── private-link.md
// │   ├── ip-filtering.png
// │   ├── ip-filtering-basic.png
// │   ├── private-link.png
// │   └── vpc-peering.png
// └── archive
//     └── 2013
//         └── meta-blogging-jekyll-setup.md
//
//  Transformations to file system structure:
//  - Put all content under a year (drop the _post / static distinction)
//  - Drop date from file name
//  - Put some content under /archive, depending on a front matter setting
//
//
// Old content front matter:
//   ---
//   layout: post
//   title: "Private Link is the IP filtering of the cloud"
//   tags: [privatelink, network]
//   tagline: Use cases for Private Link and differences in its implementation across the major Cloud Providers
//   ---
//
// New content front matter:
//   ---
//   tags: [privatelink, network]
//   date: 2022-01-30
//   ---
//   # Private Link is the IP filtering of the cloud
//
//	 Use cases for Private Link and differences in its implementation across the major Cloud Providers.
//
// Transformations to front matter
// - Drop `layout`
// - Drop `title` (but add it to content)
// - Drop `tagline` (but add it to content)
// - Add `date` (based on the prev file name)
//
// Transformations to file content:
// - Add title in the first line
// - Add tagline after the first line
// - Rewrite image urls
package main

import (
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"time"
)

func main() {
	from := flag.String("from", "", "where to migrate the content from")
	flag.Parse()

	if *from == "" {
		log.Fatalf("-from is required")
	}

	posts, err := ReadContentFromDir(filepath.Join(*from, "_posts"))
	if err != nil {
		log.Fatalf("Can't read posts from the source dir: %v\n", err)
	}
	imgs, err := ReadContentFromDir(filepath.Join(*from, "static", "img", "posts"))
	if err != nil {
		log.Fatalf("Can't read images from the source dir: %v\n", err)
	}
	cts := append(posts, imgs...)
	for i, ct := range cts {
		log.Printf("%03d: %10s, %s, %s\n", i, ct.DatePrefix, ct.Title, ct.Extension)
	}
}

func ReadContentFromDir(path string) ([]*ContentFile, error) {
	contents := make([]*ContentFile, 0)
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		datePrefix, title, ext := splitFname(filepath.Base(path))
		contents = append(contents, &ContentFile{
			Title:      title,
			DatePrefix: datePrefix,
			Extension:  ext,
			Content:    buf,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return contents, nil
}

var fnameRegexp = regexp.MustCompile(`^([0-9]{4}-[0-9]{2}(-[0-9]{2})?)-([0-9A-z-]+)`)

func splitFname(fname string) (datePrefix string, title string, ext string) {
	ext = filepath.Ext(fname)[1:]
	noExt := fname[:len(fname)-len(ext)-1]

	matches := fnameRegexp.FindStringSubmatch(noExt)
	datePrefix = matches[1]
	title = matches[3]

	return
}

// ContentFile represents a file on disk, a piece of content to be moved/transformed.
type ContentFile struct {
	Title      string
	DatePrefix string
	Extension  string
	Content    []byte
}

// Post is a text markdown post
type Post struct {
	ContentFile
	CreationDate time.Time
	Tags         []string
	Tagline      string
}

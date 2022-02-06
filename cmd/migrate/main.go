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

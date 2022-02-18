package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v2"
)

func main() {
	from := flag.String("from", "", "where to migrate the content from")
	to := flag.String("to", "", "where to migrate the content to")
	flag.Parse()

	if *from == "" {
		log.Fatalf("-from is required")
	}
	if *to == "" {
		log.Fatalf("-to is required")
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

	err = SaveContentToDir(*to, cts)
	if err != nil {
		log.Fatalf("Can't write the content to target: %v\n", err)
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
		content := &ContentFile{
			Title:      title,
			DatePrefix: datePrefix,
			Extension:  ext,
			Content:    buf,

			Post: nil,
		}
		contents = append(contents, content)

		if content.Extension != "md" {
			return nil
		}

		// post specific transformations
		post, err := processPost(content.Title, content.DatePrefix, content.Content)
		content.Post = post
		return nil
	})
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func processPost(fname string, created string, content []byte) (*Post, error) {
	var fromMatter struct {
		Title   string   `yaml:"title"`
		Tags    []string `yaml:"tags"`
		Tagline string   `yaml:"tagline"`
	}
	rest, err := frontmatter.Parse(bytes.NewReader(content), &fromMatter)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 0)

	from := fmt.Sprintf("%s/%s.html", strings.ReplaceAll(created, "-", "/"), fname)

	toMatter := &struct {
		Tags         []string `yaml:"tags"`
		Created      string   `yaml:"created"`
		RedirectFrom []string `yaml:"redirect_from"`
	}{
		Tags:         fromMatter.Tags,
		Created:      created,
		RedirectFrom: []string{from},
	}
	matterBuf, err := yaml.Marshal(toMatter)
	if err != nil {
		return nil, err
	}

	// start with front matter
	buf = append(buf, []byte(fmt.Sprintf("---\n"))...)
	buf = append(buf, matterBuf...)
	buf = append(buf, []byte(fmt.Sprintf("---\n"))...)

	// add title
	buf = append(buf, []byte(fmt.Sprintf("# %s\n", fromMatter.Title))...)

	if fromMatter.Tagline != "" {
		// add tagline
		buf = append(buf, []byte(fmt.Sprintf("\n%s\n", fromMatter.Tagline))...)
	}

	buf = append(buf, rest...)

	p := &Post{
		Content: buf,
	}
	return p, nil
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

func SaveContentToDir(path string, cts []*ContentFile) error {
	// target path must not exist, as a fail-safe not to overwrite anything
	err := os.Mkdir(path, 0755)
	if err != nil && errors.Is(err, os.ErrExist) {
		log.Printf("The target directory '%s' exists, its content maybe overwritten\n", path)
	}
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	for _, ct := range cts {
		// try creating the year directory, ignore errors because the directory might have been created already
		yearDir := filepath.Join(path, ct.Year())
		_ = os.Mkdir(yearDir, 0755)

		ctPath := filepath.Join(yearDir, ct.NewFname())
		content := ct.Content
		if ct.Post != nil {
			content = ct.Post.Content
		}
		if err := os.WriteFile(ctPath, content, 0644); err != nil {
			return err
		}
	}
	return nil
}

// ContentFile represents a file on disk, a piece of content to be moved/transformed.
type ContentFile struct {
	Title      string
	DatePrefix string
	Extension  string
	Content    []byte
	Post       *Post
}

// Year returns the year when this ContentFile was created
func (c *ContentFile) Year() string {
	return c.DatePrefix[:4]
}

func (c *ContentFile) NewFname() string {
	return fmt.Sprintf("%s.%s", c.Title, c.Extension)
}

// Post represents the extra fields for the layout: _post content
type Post struct {
	Content []byte
}

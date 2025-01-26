package crncl

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

type MetadataQuerier interface {
	Query() ([]PostMetadata, error)
}

type PostMetadata struct {
	Slug        string
	Title       string    `toml:"title"`
	Author      Author    `toml:"author"`
	Description string    `toml:"description"`
	Date        time.Time `toml:"date"`
}

type SlugReader interface {
	Read(slug string) (string, error) //TODO: should probably return an io.ReadCloser here instead of string
}

type FileReader struct {
	// Directory to find blog posts in.
	Dir string
}

func (fr FileReader) Read(slug string) (string, error) {
	slugPath := filepath.Join(fr.Dir, slug+".md")
	f, err := os.Open(slugPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (fr FileReader) Query() ([]PostMetadata, error) { // TODO: this is not most performant as it does no caching and will be run every time the index page is visited; revisit this
	postsPath := filepath.Join(fr.Dir, "*.md")
	filenames, err := filepath.Glob(postsPath)
	if err != nil {
		return nil, fmt.Errorf("querying for files: %w", err)
	}
	var posts []PostMetadata
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil { // TODO: in a prod env, there is a case for not failing whole app if one file fails; revisit this
			return nil, fmt.Errorf("opening file %s: %w", filename, err)
		}
		defer f.Close()
		var post PostMetadata
		_, err = frontmatter.Parse(f, &post)
		if err != nil {
			return nil, fmt.Errorf("parsing frontmatter for file %s: %w", filename, err)
		}
		post.Slug = strings.TrimSuffix(filepath.Base(filename), ".md")

		posts = append(posts, post)
	}
	return posts, nil
}

type PostData struct {
	Content template.HTML

	Title  string `toml:"title"`
	Author Author `toml:"author"`
}

type Author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

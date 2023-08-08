package pkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type (
	Post struct {
		Name        string
		Description string
	}

	Posts []Post

	ParsedPost struct {
		Name        string
		ID          string
		Description string
		Toc         Toc
		Content     []Section
	}

	Heading struct {
		Level int
		ID    string
		Body  string
	}

	Toc []Heading

	Content []Section

	Section struct {
		Title   string
		ID      string
		Content string
	}
)

func (s *Section) Reset() {
	s.Content = ""
	s.Title = ""
	s.ID = ""
}

func HeadingNameToID(s string) string {
  return strings.Join(strings.Split(strings.TrimSpace(s), " "), "_")
}

func SectionFrom(s []string) Section {
	sec := Section{}

	if len(s) > 0 {
    heading := strings.TrimSpace(s[0])
		if First(heading, '#') {
			title := HeadingName(s[0])
			sec.Title = title
			sec.ID = HeadingNameToID(title)
		} else {
			sec.Title = ""
			sec.ID = ""
		}
	}

	buf := []string{}
	for _, ln := range s[1:] {
		buf = append(buf, ln)
	}

	sec.Content = strings.Join(buf, "")
	return sec
}

func (posts Posts) Where(fn func(Post) bool) []Post {
	res := []Post{}
	for _, p := range posts {
		if fn(p) {
			res = append(res, p)
		}
	}
	return res
}

func (posts Posts) FirstWhere(fn func(Post) bool) (Post, error) {
	for _, p := range posts {
		if fn(p) {
			return p, nil
		}
	}
	return Post{}, errors.New("Did not find any matching element")
}

func Parse(post string, body string) (ParsedPost, error) {
	toc := ExtractToc(body)

	desc, rest, err := ExtractDescription(body)
	if err != nil {
		return ParsedPost{}, err
	}

	content := ExtractContent(rest)
	if err != nil {
		return ParsedPost{}, err
	}

	name := HeadingNameToID(post)

	return ParsedPost{
		Name:        name,
		ID:          post,
		Description: desc,
		Toc:         toc,
		Content:     content,
	}, nil
}

func ExtractDescription(s string) (string, string, error) {
	lines := strings.Split(s, "\n")
  
  parts := splitNWhere(lines, 2, func(s string) bool {
    begin := StartsWith(s, "@D") 
    if begin {
      return true
    }
    return StartsWith(s, "@Dend")
  })

  fmt.Println(len(parts))
  if len(parts) < 2 {
    return "", "", errors.New("Description not found.")
  }

  description, rest := parts[1], parts[2]
  return strings.Join(description, "\n"), strings.Join(rest, "\n"), nil
}

func ExtractContent(s string) Content {
	lines := strings.Split(s, "\n")
	content := []Section{}

	sections := splitWhere(lines, func(l string) bool {
		return First(l, '#')
	})

	for _, sec := range sections {
		section := SectionFrom(sec)
		content = append(content, section)
	}
	return content
}

func splitWhere(arr []string, fn func(string) bool) [][]string {
	parts := [][]string{}
	curr := []string{}
	for _, itm := range arr {
		if fn(itm) {
			parts = append(parts, curr)
			curr = []string{}
      curr = append(curr, itm)
			continue
		}
		curr = append(curr, itm)
	}
	return parts
}

func splitNWhere(arr []string, n int, fn func(string)bool) [][]string {
	parts := [][]string{}
	curr := []string{}
  nInternal := 0
	for i, itm := range arr {
		if fn(itm) {
			parts = append(parts, curr)
      nInternal += 1
			curr = []string{}
      if nInternal == n {
        if i + 1 < len(arr) {
          parts = append(parts, arr[i + 1:])
        }
        return parts
      }
			continue
		}
		curr = append(curr, itm)
	}
	return parts

}

func ExtractToc(s string) Toc {
	lines := strings.Split(s, "\n")
	headings := []Heading{}
	for _, line := range lines {
		if First(line, '#') {
			lvl, title := CountPrefixReturnRest(line, '#')
			headings = append(headings, Heading{
				Level: lvl,
				Body:  title,
				ID:    HeadingNameToID(title),
			})
		}
	}
	return headings
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

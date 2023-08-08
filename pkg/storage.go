package pkg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type PostStorage interface {
  Posts() Posts
  Get() Post
}

type FileStorage struct {
  path string
  posts map[string]ParsedPost
}

func NewFileStorage(path string) FileStorage {
  files, err := ioutil.ReadDir(path)
  if err != nil {
    panic("Invalid path.")
  }

  posts := make(map[string]ParsedPost)
  for _, fi := range files {
    p := path + "/" + fi.Name()
    fiBody, err := ioutil.ReadFile(p)
    if err != nil {
      fiBody = []byte("Could not read the file.")
    }
    keyName := strings.Split(fi.Name(), ".")[0]
    parsed, err := Parse(keyName, string(fiBody))
    if err != nil {
      fmt.Println(err)
      panic("Failed to parse the post.")
    }
    posts[keyName] = parsed
  }

  return FileStorage{
    path: path,
    posts: posts,
  } 
}

func (s *FileStorage) ConfirmInitialised() bool {
  return s.posts != nil
}

func (s *FileStorage) Get(file string) (ParsedPost, error) {
  if s.posts == nil {
    return ParsedPost{}, errors.New("Failed to initialise storage.")
  }

  if fi, ok := s.posts[file]; ok {
    return fi, nil
  }
  return ParsedPost{}, errors.New("Did not find file.")
}

func (s *FileStorage) Posts() []ParsedPost {
  buf := []ParsedPost{}
  for _, v := range s.posts {
    buf = append(buf, v)
  }
  return buf
}

package pkg

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func getParts(s string) ([]string, *node) {
  parts, _ := extractParts(s)
  return parts, &node{}
}

func insertDefault(n *node, parts []string) {
  n.insert(parts, func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello World")
  })
}

func TestExtractPartsRoot(t *testing.T) {
	path := "/"
	parts, err := extractParts(path)

	if err != nil {
		t.Error("Failed extracting parts")
	}

	if !reflect.DeepEqual(parts, []string{""}) {
		t.Error("Parts not equal")
	}
}

func TestExtractPartsEmpty(t *testing.T) {
	path := ""
	parts, err := extractParts(path)
	if err == nil {
		t.Error("Failed extracting parts")
	}

	if !reflect.DeepEqual(parts, []string{}) {
		t.Error("Parts not equal")
	}
}

func TestInsertOneNode(t *testing.T) {
	path := "/hello"
  parts, root := getParts(path)
  insertDefault(root, parts)

	_, exists := root.children["hello"]
	if !exists {
		t.Error("Child not found")
	}
}

func TestInsertTwoNodes(t *testing.T) {
	path := "/hello/world"
  parts, root := getParts(path)
  insertDefault(root, parts)

	hello, exists := root.children["hello"]
	if !exists {
		t.Error("Hello node not found")
	}

	_, worldExists := hello.children["world"]
	if !worldExists {
		t.Error("World node not found")
	}
}

func TestOneParam(t *testing.T) {
  path := "/hello/:name"
  parts, root := getParts(path)
  insertDefault(root, parts)
  
  hello, exists := root.children["hello"]
  if !exists {
    t.Error("Hello node not found")
  }

  _, paramExists := hello.children[":"]
  if !paramExists {
    t.Error(":param node not found")
  }
}

func TestMultipleParams(t *testing.T) {
  path := "/hello/:firstName/:secondName/:age"
  parts, root := getParts(path)
  insertDefault(root, parts)

  hello, ok := root.children["hello"]
  if !ok {
    t.Error("hello node not found")
  }

  p1, p1ok := hello.children[":"];
  if !p1ok {
    t.Error("p1 node not found")
  }

  p2, p2ok := p1.children[":"];
  if !p2ok {
    t.Error("p2 node not found")
  }

  _, p3ok := p2.children[":"];
  if !p3ok {
    t.Error("p3 node not found")
  }
}

func testSearchOneParam(t *testing.T) {
  path := "/:blog"
  parts, root := getParts(path)
  insertDefault(root, parts)

  ps, _ := extractParts(path)
  _, err := root.search(ps, map[string]string{"blog": "HelloWorld"})
  if err != nil {
    t.Error("Didnt find the parameter")
  }
}


func TestSearchNodePlusOneParam(t *testing.T) {
  path := "/hello/:name"
  parts, root := getParts(path)
  insertDefault(root, parts)

  test := "/hello/bobthebuilder"
  ps, _ := extractParts(test)
  m := map[string]string{"name": "Bob"}
  _, err := root.search(ps, m)
  if err != nil {
    t.Error("Didnt find the name parameter")
  }
}

















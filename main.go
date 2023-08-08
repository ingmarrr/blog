package main

import (
	"kinio/pkg"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var (
  storage = pkg.NewFileStorage("posts")
)

func main() {
	r := pkg.NewRouter()
	r.Use(CorsMW)
	r.Use(LogMW)
	r.Handle("/", IndexH)
	r.Handle("/post/:name", PostH)
	r.Handle("/favicon.ico", FaviconH)
  // r.Handle("/file/:name", FileH)

	http.ListenAndServe("127.0.0.1:8080", &r)
}

func LogMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		next(w, r)
	}
}

func CorsMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}
		next(w, r)
	}
}

func IndexH(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("view/index.html"))
  posts := storage.Posts()
	tmpl.Execute(w, posts)
}

func PostH(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	postName := segs[len(segs)-1]
  post, err := storage.Get(postName)
	if err != nil {
		errorTmpl := template.Must(template.ParseFiles("view/error.html"))
		errorTmpl.Execute(w, map[string]string{"Message": "Did not find this post :["})
	}
	tmpl := template.Must(template.ParseFiles("view/post.html"))
	tmpl.Execute(w, post)
}

// func FileH(w http.ResponseWriter, r *http.Request) {
//   segs := strings.Split(r.URL.Path, "/")
//   file := segs[len(segs) - 1]
//   post, err := storage.Get(file)
//   if err != nil {
//     http.NotFound(w, r)
//     return
//   }

//   w.Header().Set("Content-Type", "text/html; charset=utf-8")
//   fmt.Fprint(w, string(mdToHTML([]byte(post.String()))))
// }


func FaviconH(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "view/favicon.ico")
}








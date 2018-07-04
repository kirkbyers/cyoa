package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>
        Choose Your Own Adventure
    </title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
    <ul>
        {{range .Options}}
        <li>
            <a href="/{{.Chapter}}">{{.Text}}</a>
        </li>
        {{end}}
    </ul>
</body>
</html>
`

// NewHandler : Create http.Handler from Story map
func NewHandler(s Story, t *template.Template) http.Handler {
	if t == nil {
		t = tpl
	}
	return handler{s, t}
}

type handler struct {
	s Story
	t *template.Template
}

// ServeHTTP : main http handler
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make sure no leading or trailing space
	path := strings.TrimSpace(r.URL.Path)
	// Handle root path case
	if path == "" || path == "/" {
		path = "/intro"
	}
	// "/intro" -> "intro"
	path = path[1:]

	// if key exists in map
	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

// JSONStory : io.Reader (file) to Story map
func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// Story : map of ids -> Chapter
type Story map[string]Chapter

/*
Chapter : User chapter
Title - string: Title of chapter
Paragraphs - []string: Content of chapter
Options - []Option: options for next chapters of arc. Can be empty [].
*/
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

/*
Option : User options for next chapters in story arc
Text - string: describing option
Chapter - string: ids to link to next chapter
*/
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

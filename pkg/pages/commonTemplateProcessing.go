package pages

import (
	"fmt"
	"html/template"
	"net/http"
)

// Templates contains all of our HTML templates, cached for re-use
var Templates = template.Must(template.ParseFiles("index.html", "history.html", "data.html", "admin.html", "_header.html", "_footer.html"))

// Favicon yields a generic favicon
func Favicon(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%s\n", r.RequestURI)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=")
}

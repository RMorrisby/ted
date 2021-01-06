package pages

import "html/template"

var Templates = template.Must(template.ParseFiles("index.html", "data.html", "admin.html"))

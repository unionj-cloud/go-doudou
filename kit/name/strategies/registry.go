package strategies

import "text/template"

var Registry = make(map[string]*template.Template)

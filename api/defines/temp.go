package defines

import "sort"

type TemplateData struct {
	Title string
	Notes []Notes
}

func (t *TemplateData) Sort() {
	sort.Slice(t.Notes, func(i, j int) bool {
		return t.Notes[i].Meta.Date > t.Notes[j].Meta.Date
	})
}

package defines

import "sort"

type Note struct {
	Title string
	Ctime string
}

type TemplateData struct {
	Title string
	Notes []Note
}

func (t *TemplateData) Sort() {
	sort.Slice(t.Notes, func(i, j int) bool {
		return t.Notes[i].Ctime > t.Notes[j].Ctime
	})
}

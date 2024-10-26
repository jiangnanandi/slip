package utils

import (
	"os"
)

func WriteNote(noteDir string, note Note) error {
	file, err := os.Create(noteDir + "/" + note.Title + ".md")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(note.Body)
	if err != nil {
		return err
	}

	return nil
}

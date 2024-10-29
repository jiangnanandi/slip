package utils

import (
	"os"
	"slip/api/types"
)

func WriteNote(noteDir string, note types.Notes) error {
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

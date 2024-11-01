package utils

import (
	"os"
	"slip/api/defines"
)

func SaveNote(note defines.Notes) error {
	file, err := os.Create(note.Dir + "/" + note.Title + ".md")
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

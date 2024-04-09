package service

import "slip/utils"

func SaveNote(note utils.Note) error {
	// 将 Note 内容写入到指定目录
	noteDir := utils.GetNoteDir()
	if err := utils.CreateDir(); err != nil {
		return err
	}
	// 写入文件
	if err := utils.WriteNote(noteDir, note); err != nil {
		return err
	}
	return nil
}

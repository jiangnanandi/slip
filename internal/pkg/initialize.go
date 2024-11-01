package initialize

import (
    "os"
    "slip/internal/config"
)


func InitializeDirectories() error {
    dirs := []string{
        config.AppConfig.Notes.PublishedDir,
        config.AppConfig.Notes.DraftDir,
        config.AppConfig.Notes.ArchivedDir,
        config.AppConfig.Notes.DeletedDir,
        config.AppConfig.Notes.PrivateDir,
    }

    for _, dir := range dirs {
        if err := os.MkdirAll(dir, 0777); err != nil {
            return err
        }
    }

    return nil
}
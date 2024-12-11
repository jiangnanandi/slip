package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slip/api/defines"
	"slip/internal/config"
	"time"
)

// RetryTask 重试任务结构
type RetryTask struct {
	Note      defines.Notes `json:"note"`
	Attempts  int           `json:"attempts"`
	CreatedAt time.Time     `json:"created_at"`
}

// SaveNote 尝试保存笔记，如果失败则加入重试队列
func SaveNote(note defines.Notes) error {
	if err := saveNoteOnce(note); err != nil {
		// 保存失败，加入重试队列
		task := RetryTask{
			Note:      note,
			Attempts:  0,
			CreatedAt: time.Now(),
		}
		return addToRetryQueue(task)
	}
	return nil
}

// saveNoteOnce 执行单次保存笔记的操作
func saveNoteOnce(note defines.Notes) error {
    // 确保目标目录存在
    if err := os.MkdirAll(note.Dir, 0755); err != nil {
        return err
    }

    filename := note.Title + ".md"

    // 删除所有目录中的同名笔记
    dirs := []string{
        config.AppConfig.Notes.PublishedDir,
        config.AppConfig.Notes.DraftDir,
        config.AppConfig.Notes.ArchivedDir,
        config.AppConfig.Notes.DeletedDir,
        config.AppConfig.Notes.PrivateDir,
    }

    for _, dir := range dirs {
        if dir == note.Dir {
            continue // 跳过目标目录
        }
        oldPath := filepath.Join(dir, filename)
        if _, err := os.Stat(oldPath); err == nil {
            if err := os.Remove(oldPath); err != nil {
                return err
            }
        }
    }

    // 创建临时文件
    tempFile := filepath.Join(note.Dir, filename+".tmp")
    file, err := os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
    if err != nil {
        return err
    }

    // 使用 defer 确保文件最终会被关闭
    defer func() {
        file.Close()
        // 如果临时文件还存在，清理它
        os.Remove(tempFile)
    }()

    // 写入内容
    if _, err = file.WriteString(note.Body); err != nil {
        return err
    }

    // 确保所有数据都写入磁盘
    if err = file.Sync(); err != nil {
        return err
    }

    // 关闭文件，准备重命名
    if err = file.Close(); err != nil {
        return err
    }

    // 原子性地重命名临时文件为目标文件
    return os.Rename(tempFile, filepath.Join(note.Dir, filename))
}

// 获取重试间隔时间
func getRetryInterval(attempts int) time.Duration {
	switch {
	case attempts <= 3:
		return 3 * time.Second
	case attempts <= 10:
		return 1 * time.Second
	case attempts <= 20:
		return 500 * time.Millisecond
	default:
		return 200 * time.Millisecond
	}
}

// 添加到重试队列
func addToRetryQueue(task RetryTask) error {
	queueDir := filepath.Join(config.AppConfig.DataDir, "retry_queue")
	if err := os.MkdirAll(queueDir, 0755); err != nil {
		return err
	}

	// 使用时间戳和随机数生成唯一文件名
	filename := fmt.Sprintf("%d_%d.json", time.Now().UnixNano(), rand.Int63())
	path := filepath.Join(queueDir, filename)

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ProcessRetryQueue 处理重试队列的后台任务
func ProcessRetryQueue(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tasks := loadRetryTasks()
			for _, task := range tasks {
				interval := getRetryInterval(task.Attempts)
				if time.Since(task.CreatedAt) < interval {
					continue // 还未到重试时间
				}

				if err := saveNoteOnce(task.Note); err != nil {
					// 先删除旧文件（包含旧的 Attempts 值）
					removeTaskFile(task)
					// 更新重试次数并保存回队列
					task.Attempts++
					addToRetryQueue(task)
				} else {
					// 成功后删除任务文件
					removeTaskFile(task)
				}
			}
		}
	}
}

// loadRetryTasks 加载所有重试任务
func loadRetryTasks() []RetryTask {
	queueDir := filepath.Join(config.AppConfig.DataDir, "retry_queue")
	files, err := os.ReadDir(queueDir)
	if err != nil {
		return nil
	}

	var tasks []RetryTask
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(queueDir, file.Name()))
		if err != nil {
			continue
		}

		var task RetryTask
		if err := json.Unmarshal(data, &task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks
}

// removeTaskFile 删除重试任务文件
func removeTaskFile(task RetryTask) error {
    queueDir := filepath.Join(config.AppConfig.DataDir, "retry_queue")
    files, err := os.ReadDir(queueDir)
    if err != nil {
        return err
    }

    for _, file := range files {
        if filepath.Ext(file.Name()) != ".json" {
            continue
        }

        path := filepath.Join(queueDir, file.Name())
        // 先尝试读取文件
        data, err := os.ReadFile(path)
        if err != nil {
            continue
        }

        var currentTask RetryTask
        if err := json.Unmarshal(data, &currentTask); err != nil {
            continue
        }

        // 使用更可靠的比较方式
        if currentTask.CreatedAt.Equal(task.CreatedAt) && 
           currentTask.Note.Title == task.Note.Title &&
           currentTask.Note.Dir == task.Note.Dir {
            return os.Remove(path)
        }
    }

    return fmt.Errorf("task file not found")
}
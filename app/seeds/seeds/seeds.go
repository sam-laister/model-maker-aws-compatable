package seeds

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Soup666/diss-api/utils"
)

func MakeBackup() error {
	srcDirs := [2]string{"./uploads", "./objects"}
	destDir := "./backup/"

	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	for _, srcDir := range srcDirs {
		os.Rename(srcDir, filepath.Join(destDir, filepath.Base(srcDir)))
	}

	return nil
}

func CopyRawModel(taskID uint) error {
	srcDir := "./seeds/backup/models"
	destDir := fmt.Sprintf("./objects/%d", taskID)

	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(destDir, relPath)
			err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			if err != nil {
				return err
			}

			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func CopyRawImages(taskID uint) error {
	srcDir := "./seeds/backup/images"
	destDir := fmt.Sprintf("./uploads/%d", taskID)

	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(srcDir, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(destDir, relPath)
			err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			if err != nil {
				return err
			}

			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func CopyFilesFrom7z() error {
	util := utils.NewFileUtil()

	currentDir, _ := os.Getwd()
	srcFile := filepath.Join(currentDir, "seeds", "backup", "backup-min.7z")

	if err := util.Extract7z(srcFile, "./"); err != nil {
		fmt.Printf("Extraction failed: %v\n", err)
		return nil
	}

	fmt.Println("Extraction successful.")

	return nil
}

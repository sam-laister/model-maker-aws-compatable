package seeds

import (
	"github.com/Soup666/diss-api/model"
	"gorm.io/gorm"
)

func CreateFile(db *gorm.DB, taskID uint, url string, filename string, filetype string) error {
	return db.Create(&model.AppFile{TaskId: taskID, Url: url, Filename: filename, FileType: filetype}).Error
}

var filenames = []string{
	"20150131_141819.jpg",
	"20150131_141828.jpg",
	"20150131_141832.jpg",
	"20150131_141837.jpg",
	"20150131_141841.jpg",
	"20150131_141845.jpg",
	"20150131_141849.jpg",
	"20150131_141853.jpg",
	"20150131_141857.jpg",
	"20150131_141902.jpg",
	"20150131_141908.jpg",
	"20150131_141913.jpg",
	"20150131_141921.jpg",
	"20150131_141927.jpg",
	"20150131_141931.jpg",
	"20150131_141935.jpg",
	"20150131_141938.jpg",
	"20150131_141942.jpg",
	"20150131_141947.jpg",
	"20150131_141951.jpg",
	"20150131_141955.jpg",
	"20150131_142000.jpg",
	"20150131_142004.jpg",
}

func CreateDummyFiles(db *gorm.DB) error {
	for _, filename := range filenames {
		if err := CreateFile(db, 1, "/uploads/1/"+filename, filename, "upload"); err != nil {
			return err
		}
	}
	return nil
}

package validate

import (
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/data"
)

func EpisodeWebURL(db *gorm.DB, start int, limit int) (episodes []data.Episode, err error) {
	dbQ := db
	if start > 0 {
		dbQ = db.Where("id <= ?", start)
	}
	err = dbQ.Where("banned = 0 AND src_type = ?", 11).
		Order("id desc").Limit(limit).Find(&episodes).Error
	return
}

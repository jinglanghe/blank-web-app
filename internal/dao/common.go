package dao

import (
	"gorm.io/gorm"
	"strings"
)

func sort(db *gorm.DB, sortStr string) *gorm.DB {
	if len(sortStr) != 0 {
		sortS := strings.Split(strings.TrimSpace(sortStr), ",")
		for _, s := range sortS {
			sortTmp := strings.Split(strings.TrimSpace(s), "|")
			if len(sortTmp) != 2 {
				continue
			}
			if sortTmp[1] != "desc" && sortTmp[1] != "asc" {
				continue
			}
			key := sortTmp[0]
			if key == "createdAt" {
				key = "created_at"
			}
			db = db.Order(key + " " + sortTmp[1])
		}
	} else {
		db = db.Order("created_at DESC")
	}

	return db
}

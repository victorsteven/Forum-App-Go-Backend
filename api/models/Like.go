package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Like struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID  uint32    `gorm:"not null" json:"user_id"`
	PostID  uint64    `gorm:"not null" json:"post_id"`
	//Post 	Post
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


func (l *Like) SaveLike(db *gorm.DB) (*Like, error) {
	var err error
	err = db.Debug().Model(&Like{}).Create(&l).Error
	if err != nil {
		return &Like{}, err
	}
	return l, nil
}

func (l *Like) DeleteLike(db *gorm.DB, pid uint64, uid uint32) (*Like, error) {

	var err error
	var deletedLike *Like

	err = db.Debug().Model(Like{}).Where("post_id = ? and user_id = ?", pid, uid).Take(&l).Error
	if err != nil {
		return &Like{}, err
	} else {
		//If the like exist, save it in deleted like and delete it
		deletedLike = l
		db = db.Debug().Model(&Like{}).Where("user_id = ?", l.UserID).Take(&Like{}).Delete(&Like{})
		if db.Error != nil {
			return &Like{}, db.Error
		}
	}
	return deletedLike, nil
}


func (l *Like) GetLikesInfo(db *gorm.DB, pid uint64) (*[]Like, error)  {
	//err := db.Debug().Model(&Like{}).Where("post_id = ?", pid).Take(&Like{}).Error
	//if err != nil {
	//	return &Like{}, err
	//}
	//return l, nil

	likes := []Like{}
	err := db.Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Error
	if err != nil {
		return &[]Like{}, err
	}
	return &likes, err

	//var result uint64
	//db.Debug().Model(&Like{}).Where("post_id = ?", pid).Count(&result)

	//check if the authenticated user liked the post:
}

func (l *Like) authUserLike(db *gorm.DB, uid uint32, pid uint64) (*Like, error) {
	err := db.Debug().Model(&Like{}).Where("post_id = ? and user_id = ?", pid, uid).Error
	if err != nil {
		return &Like{}, err
	}
	return l, nil
}
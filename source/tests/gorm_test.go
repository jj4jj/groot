package tests

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"groot/sfw/util"
	"testing"
	"time"
	sfw_db "groot/sfw/db"
)


//testing gorm sql injection

type DbDemo struct {
	gorm.Model
	Str 	string
	Byt 	[]byte	`gorm:"size:16384"`
	Boo 	bool
	In64 	int64
	Tm 		time.Time
	tm2 	*time.Time
}

func init(){
	util.CheckError(sfw_db.AddDb("test", "sqlite3", "./test.db",
		0,0),"")
}

func TestCreate(t *testing.T) {
	db := sfw_db.GetDb("test")
	db.AutoMigrate(&DbDemo{})
	db.Create(&DbDemo{
		Str:"1",
	})
	db.Create(&DbDemo{
		Str:"2",
	})
	db.Create(&DbDemo{
		Str:"3",
	})
	db.Create(&DbDemo{
		Str:"4",
	})
}
func TestSqlInject(t *testing.T) {
	db := sfw_db.GetDb("test")
	dm := DbDemo{}
	var dblist []DbDemo
	db.Debug().Model(&dm).Find(&dblist)
	fmt.Println(dblist)

	q := "' OR '1' != '"
	db.Debug().Model(&dm).Find(&dblist, " Str = ?", q)
	fmt.Println(dblist)

	db.Debug().Delete(&dm, "Str = ?", q)
}











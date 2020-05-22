package sfw_db

import (
	"errors"
	"github.com/lexkong/log"
	"reflect"
	"strings"
	"unicode"
)

type UpdateFieldsMap map[string]interface{}

func GetModelTableNameBase(val interface{}) string {
	t := reflect.TypeOf(val)
	kind := t.Kind()
	if kind != reflect.Struct && kind != reflect.Ptr {
		return ""
	}
	//pkg.DbXxName => db_xx_name
	n := t.String()
	if n == "" {
		return ""
	}

	as := strings.Split(n, ".")

	if len(as) != 2 {
		return ""
	}
	var result = []rune("tb_")
	rn := []rune(as[1])
	for i,r := range rn {
		if i > 0 && unicode.IsLower(rn[i-1]) && unicode.IsUpper(rn[i]) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func AutoMigrateTable(DbNamBase string,values ...interface{}) error {
	return 	AutoMigrateSplitDbAndTables(DbNamBase, 0,0, values...)
}

func GetModelTableName(value interface{}) string {
	return GetModelTableNameBase(value)
}

func GetModelTableNameWithSplitFactor(SplitDbIdx,SplitTbIdx uint32,value interface{}) string {
	return GetTbNameByIdx(GetModelTableNameBase(value), SplitDbIdx, SplitTbIdx)
}

func AutoMigrateSplitDbAndTables(DbNameBase string, SplitDbNum, SplitTbNum uint32, values ...interface{}) error {

	itrDbNum := SplitDbNum
	if SplitDbNum == 0 {
		itrDbNum = 1
	}
	itrTbNum := SplitTbNum
	if SplitTbNum == 0 {
		itrTbNum = 1
	}
	for i:= uint32(0); i < itrDbNum; i++ {
		dbname := GetDbNameByIdx(DbNameBase, i)
		if SplitDbNum == 0 {
			dbname = DbNameBase
		}
		db := GetDb(dbname)
		if db == nil {
			return errors.New("db name not exist")
		}
		for j := uint32(0); j < itrTbNum; j++ {
			for k := 0; k < len(values); k++ {
				val := values[k]
				tbname_base := GetModelTableNameBase(val)
				if tbname_base == "" {
					log.Errorf(nil,"get table name base error")
					return errors.New("error table name base get")
				}
				tbname := GetTbNameByIdx(tbname_base, i, j)
				if SplitTbNum == 0 {
					tbname = tbname_base
				}
				log.Debugf("auto migrate db table:%s.%s ...", dbname, tbname)
				if e:=db.Table(tbname).AutoMigrate(val).Error; e!=nil {
					log.Errorf(e, "auto migrate dbname:%s tbname:%s (i:%d/%d,j:%d/%d) error",
						dbname, tbname, i, SplitDbNum, j, SplitTbNum)
					return e
				}
			}
		}
	}
	return nil
}


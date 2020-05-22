package sfw_db

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

type GetKeyDbIndexFunc func(key []byte) uint32
type GetKeyTbIndexFunc func(key []byte) uint32

func GetDbKeyHashCode(key []byte) uint32 {
	w := md5.New()
	b := w.Sum(key)
	return binary.LittleEndian.Uint32(b)
}

func AutoScaleCreateDbKeyFunc(maxDbNum,maxTbNum uint32) (GetKeyDbIndexFunc,GetKeyTbIndexFunc){
	f1 := func(key []byte) uint32 {
		code := GetDbKeyHashCode(key)
		return (code/maxTbNum) % maxDbNum
	}
	f2 := func(key []byte) uint32 {
		code := GetDbKeyHashCode(key)
		return code % maxTbNum
	}
	return f1, f2
}

func GetDbNameByIdx(dbnameBase string, dbidx uint32) string {
	return fmt.Sprintf("%s_%d", dbnameBase, dbidx)
}

func GetTbNameByIdx(tbnameBase string, dbidx,tbidx uint32) string {
	return fmt.Sprintf("%s_%d_%d", tbnameBase, dbidx, tbidx)
}

func GetDbNameByKey(key []byte, dbsplitNum uint32, dbnameBase string) string {
	code := GetDbKeyHashCode(key)
	return GetDbNameByIdx(dbnameBase, code%dbsplitNum)
}

func GetTbNameByKey(key []byte, dbsplitNum, tbsplitNum uint32, tbnameBase string) string {
	code := GetDbKeyHashCode(key)
	dbidx := (code/ tbsplitNum) % dbsplitNum
	return GetTbNameByIdx(tbnameBase, dbidx, code%dbsplitNum)
}



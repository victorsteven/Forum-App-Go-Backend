package fileformat

import (
	"path"
	"github.com/twinj/uuid"
	"strings"
	"time"
)


func UniqueFormat(fn string) string {
	//path.Ext() get the extension of the file
	fileName :=  strings.TrimSuffix(fn, path.Ext(fn))
	extension := path.Ext(fn)
	t := time.Now()
	u := uuid.NewV4()
	 newFileName := fileName + "--" +  t.String() + "--" +  u.String() + extension

	 return newFileName


}
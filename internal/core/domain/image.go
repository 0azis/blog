package domain

import "mime/multipart"

type Image struct {
	File *multipart.FileHeader
}

func (i Image) IsValid() bool {
	if i.File.Size < 50000 {
		return false
	}
	return true
}

package model

type UploadImageInput struct {
	FileData []byte
	FileName string
	Folder   string
}

type UploadImageRespose struct {
	Fileid    string
	Url       string
	Thumbnail string
}

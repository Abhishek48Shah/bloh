package blog

type BlogServices struct {
	storage *BlogStorage
}

func LoadServices(storage *BlogStorage) *BlogServices {
	return &BlogServices{storage: storage}
}

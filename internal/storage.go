package internal

type DB interface {
	GetShortURL(url string) (string, error)
	GetLongURL(url string) (string, error)
	AddURL(url string) (string, error)
}

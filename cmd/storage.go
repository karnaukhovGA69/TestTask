package main

type DB interface {
	GetShortURl(url string) (string, error)
	GetLongURL(url string) (string, error)
	AddURL(url string) (string, error)
}

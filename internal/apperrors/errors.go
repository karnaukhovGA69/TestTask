package apperrors

import "errors"

var (
	ErrBadJSON                  = errors.New("не получилось распарсить URL")
	ErrEmptyURL                 = errors.New("пустой URL")
	ErrNotFound                 = errors.New("ссылка не найдена")
	ErrShortURLGeneration       = errors.New("не получилось создать короткий URL")
	ErrUniqueShortURLGeneration = errors.New("не получилось создать уникальный короткий URL")
	ErrMissingConfig            = errors.New("не хватает переменной окружения")
	ErrUnknownStorage           = errors.New("unknown database type")
)

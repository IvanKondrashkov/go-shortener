package storage

import "errors"

// Пакет storage содержит определения ошибок хранилища
var (
	// ErrConflict - возникает при попытке создать дублирующую сущность
	ErrConflict = errors.New("entity conflict")
	// ErrBatchIsEmpty - возникает при обработке пустого пакета данных
	ErrBatchIsEmpty = errors.New("batch is empty")
	// ErrURLNotValid - возникает при передаче невалидного URL
	ErrURLNotValid = errors.New("url is invalidate")
	// ErrNotFound - возникает при запросе несуществующей сущности
	ErrNotFound = errors.New("entity not found")
	// ErrDeleteAccepted - возникает при успешном принятии запроса на удаление
	ErrDeleteAccepted = errors.New("entity accepted delete")
)

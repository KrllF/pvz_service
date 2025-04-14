package errs

import "errors"

var (
	// ErrPackTypeNotFound - тип упаковки не существует
	ErrPackTypeNotFound = errors.New("тип упаковки не существует")
	// ErrStatusNotFound - статус заказа не существует
	ErrStatusNotFound = errors.New("статус заказа не существует")
	// ErrExtraPackNotFound - дополнительная упаковка не найдена
	ErrExtraPackNotFound = errors.New("дополнительная упаковка не найдена")
	// ErrOrderNotFound - заказ не найден
	ErrOrderNotFound = errors.New("заказ не найден")
	// ErrUserNotFound - пользователь не найден
	ErrUserNotFound = errors.New("пользователь не найден")
	// ErrFileNotFound - файл не найден
	ErrFileNotFound = errors.New("файл не найден")

	// ErrInvalidData - невалидные данные
	ErrInvalidData = errors.New("невалидные данные")
	// ErrInvalidSortField - недопустимое поле для сортировки
	ErrInvalidSortField = errors.New("недопустимое поле для сортировки")
	// ErrInvalidQueryParams - некорректные параметры запроса
	ErrInvalidQueryParams = errors.New("некорректные параметры запроса")
	// ErrInvalidOrderStatus - недопустимый статус заказа
	ErrInvalidOrderStatus = errors.New("недопустимый статус заказа")

	// ErrTimeNotExpired - время для удаления заказа еще не истекло
	ErrTimeNotExpired = errors.New("время для удаления заказа еще не истекло")
	// ErrShelfLifeExpired - срок хранения заказа истёк
	ErrShelfLifeExpired = errors.New("срок хранения заказа истёк")

	// ErrInternalServerError - внутренняя ошибка сервера
	ErrInternalServerError = errors.New("внутренняя ошибка сервера")
	// ErrFileReadError - ошибка при чтении файла
	ErrFileReadError = errors.New("ошибка при чтении файла")
)

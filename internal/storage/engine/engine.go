package engine

import (
	"errors"
	"go.uber.org/zap"
)

var ErrKeyNotFound = errors.New("key not found")

type Engine struct {
	storage map[string]string
	logger  *zap.Logger
}

func NewEngine(logger *zap.Logger) *Engine {
	return &Engine{
		storage: make(map[string]string),
		logger:  logger,
	}
}

func (e *Engine) Set(key string, value string) error {
	e.storage[key] = value

	return nil
}

func (e *Engine) Get(key string) (string, error) {
	value, ok := e.storage[key]
	if !ok {
		e.logger.Info("Get: entry not found",
			zap.String("key", key),
		)
		return "", ErrKeyNotFound
	}

	return value, nil
}

func (e *Engine) Delete(key string) error {
	_, ok := e.storage[key]
	if !ok {
		e.logger.Info("Delete: entry not found",
			zap.String("key", key),
		)
	}

	delete(e.storage, key)
	return nil
}

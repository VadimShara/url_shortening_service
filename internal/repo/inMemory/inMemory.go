package inMemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	errs "github.com/VadimShara/url_shortening_service/pkg/errs"
)

type DB struct {
	mu            sync.RWMutex
	urlStorage    map[string]string // stores url -> alias
	aliasStorage  map[string]string // stores alias -> url
	clearInterval time.Duration
}

func NewDB() *DB {
	db := &DB{
		urlStorage:    make(map[string]string),
		aliasStorage:  make(map[string]string),
		clearInterval: time.Minute,
	}

	return db
}

func (d *DB) SaveUrl(ctx context.Context, urlToSave, alias string) (string, error) {
	const op = "repo.inMemory.SaveUrl"

	d.mu.Lock()
	defer d.mu.Unlock()

	if existingAlias, exists := d.urlStorage[urlToSave]; exists {
		return existingAlias, fmt.Errorf("%s: %w", op, errs.ErrUrlExists)
	}

	if _, exists := d.aliasStorage[alias]; exists {
		return "", fmt.Errorf("%s: %w", op, errs.ErrAliasExists)
	}

	d.urlStorage[urlToSave] = alias
	d.aliasStorage[alias] = urlToSave

	return alias, nil
}

func (d *DB) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "repo.inMemory.GetUrl"

	d.mu.RLock()
	defer d.mu.RUnlock()

	url, exists := d.aliasStorage[alias]
	if !exists {
		return "", fmt.Errorf("%s: %w", op, errs.ErrUrlNotFound)
	}

	return url, nil
}

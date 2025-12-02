package models

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type TodoStatus string

const (
	StatusPending    TodoStatus = "pending"
	StatusInProgress TodoStatus = "in-progress"
	StatusDone       TodoStatus = "done"
)

func (s TodoStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	}
	return false
}

func (s TodoStatus) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, errors.New("invalid status")
	}
	return string(s), nil
}

func (s *TodoStatus) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("invalid type for status")
	}
	v = strings.ToLower(v)
	tmp := TodoStatus(v)
	if !tmp.IsValid() {
		return errors.New("invalid status")
	}
	*s = tmp

	return nil
}

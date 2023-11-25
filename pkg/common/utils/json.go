package utils

import (
	"encoding/json"
	"time"
)

type JSONString struct {
	Value string
	Valid bool
	Set   bool
}

func (s *JSONString) UnmarshalJSON(data []byte) error {
	s.Set = true

	if string(data) == "null" {
		s.Valid = false
		return nil
	}

	var temp string

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	s.Value = temp
	s.Valid = true
	return nil
}

type JSONStringSlice struct {
	Value []string
	Valid bool
	Set   bool
}

func (ss *JSONStringSlice) UnmarshalJSON(data []byte) error {
	ss.Set = true

	if string(data) == "null" {
		ss.Valid = false
		return nil
	}

	var temp []string

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	ss.Value = temp
	ss.Valid = true
	return nil
}

type JSONTime struct {
	Value time.Time
	Valid bool
	Set   bool
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	t.Set = true

	if string(data) == "null" {
		t.Valid = false
		return nil
	}

	var temp time.Time

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	t.Value = temp
	t.Valid = true
	return nil
}

type JSONBool struct {
	Value bool
	Valid bool
	Set   bool
}

func (b *JSONBool) UnmarshalJSON(data []byte) error {
	b.Set = true

	if string(data) == "null" {
		b.Valid = false
		return nil
	}

	var temp bool

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	b.Value = temp
	b.Valid = true
	return nil
}

package immutable

import (
	"encoding/json"
	"testing"
)

func TestSome(t *testing.T) {
	opt := Some(1)
	if !opt.HasValue() {
		t.Errorf("expected Some to return an Option with a value")
	}
	if opt.Value() != 1 {
		t.Errorf("expected Some to return an Option with a value of 1")
	}
}

func TestNone(t *testing.T) {
	opt := None[int]()
	if opt.HasValue() {
		t.Errorf("expected None to return an Option with no value")
	}
}

func TestOptionMarshal(t *testing.T) {
	opt := Some(1)
	b, err := json.Marshal(opt)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if string(b) != "1" {
		t.Errorf("expected 1, got %s", b)
	}
}

func TestOptionMarshalNone(t *testing.T) {
	opt := None[int]()
	b, err := json.Marshal(opt)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if string(b) != "null" {
		t.Errorf("expected null, got %s", b)
	}
}

func TestOptionUnmarshal(t *testing.T) {
	var opt Option[int]
	err := json.Unmarshal([]byte("1"), &opt)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !opt.HasValue() {
		t.Errorf("expected Some to return an Option with a value")
	}
	if opt.Value() != 1 {
		t.Errorf("expected Some to return an Option with a value of 1")
	}
}

func TestOptionUnmarshalNone(t *testing.T) {
	var opt Option[int]
	err := json.Unmarshal([]byte("null"), &opt)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if opt.HasValue() {
		t.Errorf("expected None to return an Option with no value")
	}
}

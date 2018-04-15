package models

import (
	"reflect"
	"testing"
	"time"
)

func TestGetAstronauts(t *testing.T) {
	a := NewAdvert()

	if reflect.TypeOf(*a).Name() != "Advert" {
		t.Errorf("Type of struct is incorrect, got: %s, want: %s.", reflect.TypeOf(a).Name(), "Advert")
	}

	if time.Now().Unix()-a.Created > 1 {
		t.Error("Created timestamp is not set")
	}

	if a.Status != statusNew {
		t.Error("Status is wrong, got: %s, want: %s", a.Status, statusNew)
	}
}

func TestSetCreatedAt(t *testing.T) {
	a := new(Advert)
	a.SetCreatedAt()
	if time.Now().Unix()-a.Created > 1 {
		t.Error("Created timestamp is not set")
	}
}

func TestSetStatusInitial(t *testing.T) {
	a := new(Advert)
	a.SetStatusInitial()
	if a.Status != statusNew {
		t.Error("Status is wrong, got: %s, want: %s", a.Status, statusNew)
	}
}

func TestGetCreated(t *testing.T) {
	a := new(Advert)
	a.SetCreatedAt()
	ti := a.GetCreated()
	if reflect.TypeOf(ti).Name() != "Time" {
		t.Errorf("Type of created field is wrong, got: %s, want: %s.", reflect.TypeOf(ti).Name(), "Time")
	}

	if time.Now().Unix()-ti.Unix() > 1 {
		t.Error("Time stored in created field is wrong, got: %s, want: %s", ti.Unix(), time.Now().Unix())
	}
}

package main

import "testing"

func TestNewModelStartsAtChoiceScreen(t *testing.T) {
	m := NewModel()
	if m.state != stateChoice {
		t.Fatalf("expected stateChoice, got %v", m.state)
	}
}

func TestLoadItemsUsesFrontendCatalog(t *testing.T) {
	m := NewModel()
	m.catalog = FrontendCategories
	m.loadItems()

	if len(m.items) == 0 {
		t.Fatal("expected frontend items to be loaded")
	}
	if m.items[0].name == "" {
		t.Fatal("expected first item to have a name")
	}
}

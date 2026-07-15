package main

import "testing"

func TestStartSelectionUsesFrontendCatalog(t *testing.T) {
	m := NewModel()
	m.startSelection("frontend")

	if m.state != stateSelecting {
		t.Fatalf("expected stateSelecting, got %v", m.state)
	}
	if m.choice != "frontend" {
		t.Fatalf("expected frontend choice, got %q", m.choice)
	}
	if len(m.catalog) != len(FrontendCategories) {
		t.Fatalf("expected %d categories, got %d", len(FrontendCategories), len(m.catalog))
	}
}

func TestBuildRowsIncludesHeadersAndLibraries(t *testing.T) {
	rows := buildRows(FrontendCategories)
	if len(rows) == 0 {
		t.Fatal("expected rows to be built")
	}

	foundHeader := false
	foundLib := false
	for _, r := range rows {
		if r.kind == rowHeader && r.header == "Routing" {
			foundHeader = true
		}
		if r.kind == rowLib && r.catIdx == 0 {
			foundLib = true
		}
	}

	if !foundHeader {
		t.Fatal("expected Routing header row")
	}
	if !foundLib {
		t.Fatal("expected first category library row")
	}
}

package models

import "testing"

func TestWorkingPartyCategoryValidate(t *testing.T) {
	category := &WorkingPartyCategory{
		Name:         "Clinical Guidelines",
		Slug:         "clinical-guidelines",
		Description:  "Working party documents",
		DisplayOrder: 1,
		DropboxPath:  "WORKING_PARTIES/Clinical%20Guidelines",
	}

	if err := category.Validate(); err != nil {
		t.Fatalf("expected valid category, got error: %v", err)
	}

	if got := category.GetDropboxPath(); got != "WORKING_PARTIES/Clinical%20Guidelines" {
		t.Fatalf("GetDropboxPath() = %q, want stored path", got)
	}
}

func TestWorkingPartyCategoryGetDropboxPathFallback(t *testing.T) {
	category := &WorkingPartyCategory{Name: "Ethics Board"}
	if got := category.GetDropboxPath(); got != "WORKING_PARTIES/Ethics Board" {
		t.Fatalf("GetDropboxPath() = %q, want WORKING_PARTIES/Ethics Board", got)
	}
}

func TestCreateWorkingPartyCategoryRequestValidate(t *testing.T) {
	req := &CreateWorkingPartyCategoryRequest{
		Name:         "Research Committee",
		Description:  "Committee documents",
		DisplayOrder: 2,
	}

	if err := req.Validate(); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}

func TestUpdateWorkingPartyCategoryRequestValidate(t *testing.T) {
	name := "Updated Committee"
	req := &UpdateWorkingPartyCategoryRequest{Name: &name}

	if err := req.Validate(); err != nil {
		t.Fatalf("expected valid update request, got error: %v", err)
	}
}

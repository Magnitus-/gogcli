package manifest

import (
	"testing"
)

func TestManifestGameExtraHasOneOfTypeTerms(t *testing.T) {
	extra := ManifestGameExtra{
		Url:           "/dontknowdontcare",
		Title:         "extra",
		Name:          "extra",
		Type:          "ost",
		Info:          1,
		EstimatedSize: "1kb",
		VerifiedSize:  1000,
		Checksum:      "dfsdfsdfwe",
	}

	if extra.HasOneOfTypeTerms([]string{"manual", "wallpaper"}) {
		t.Errorf("Should not indicate it has a type it doesn't have")
	}

	if !extra.HasOneOfTypeTerms([]string{"manual", "ost"}) {
		t.Errorf("Should indicate it has a type it has")
	}
}

func TestManifestGameExtraIsEquivalentTo(t *testing.T) {
	extra := ManifestGameExtra{
		Url:           "/dontknowdontcare",
		Title:         "extra",
		Name:          "extra",
		Type:          "ost",
		Info:          1,
		EstimatedSize: "1kb",
		VerifiedSize:  1000,
		Checksum:      "abcdefgh",
	}

	otherExtra := ManifestGameExtra{
		Url:           "/dontknowdontcare",
		Title:         "extra",
		Name:          "extra",
		Type:          "manual",
		Info:          2,
		EstimatedSize: "2kb",
		VerifiedSize:  1000,
		Checksum:      "abcdefgh",
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumValidation, false) {
		t.Errorf("Extras who url, title, name, verified size and checksum match should be equivalent when empty checksums are not tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras who url, title, name, verified size and checksum match should be equivalent when empty checksums are tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumNoValidation, false) {
		t.Errorf("Extras who url, title, name, verified size and checksum match should be equivalent when checksums are not validated")
	}

	otherExtra.Checksum = "idonotcheck"
	if extra.IsEquivalentTo(&otherExtra, ChecksumValidation, false) {
		t.Errorf("Extras whose checksum doesn't match should not be equivalent when empty checksums are not tolerated")
	}

	if extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras whose checksum doesn't match should not be equivalent when empty checksums are tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumNoValidation, false) {
		t.Errorf("Extras whose checksum doesn't match should be equivalent when checksums are not validated")
	}

	otherExtra.Checksum = ""
	if extra.IsEquivalentTo(&otherExtra, ChecksumValidation, false) {
		t.Errorf("Extras whose checksum is empty should not be equivalent when empty checksums are not tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras whose checksum is empty should be equivalent when empty checksums are tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumNoValidation, false) {
		t.Errorf("Extras whose checksum is empty should be equivalent when checksums are not validated")
	}

	extra.Checksum = ""
	otherExtra.Checksum = "icheck"
	if extra.IsEquivalentTo(&otherExtra, ChecksumValidation, false) {
		t.Errorf("Extras whose checksum is empty should not be equivalent when empty checksums are not tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras whose checksum is empty should be equivalent when empty checksums are tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumNoValidation, false) {
		t.Errorf("Extras whose checksum is empty should be equivalent when checksums are not validated")
	}

	extra.Checksum = ""
	otherExtra.Checksum = ""
	if extra.IsEquivalentTo(&otherExtra, ChecksumValidation, false) {
		t.Errorf("Extras whose checksum is empty should not be equivalent when empty checksums are not tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras whose checksum is empty should be equivalent when empty checksums are tolerated")
	}

	if !extra.IsEquivalentTo(&otherExtra, ChecksumNoValidation, false) {
		t.Errorf("Extras whose checksum is empty should be equivalent when checksums are not validated")
	}

	otherExtra.Url = "/dadada"
	if extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras whose url differ should not be equivalent")
	}

	otherExtra.Url = "/dontknowdontcare"
	otherExtra.Title = "wrong"
	if extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras whose title differ should not be equivalent")
	}

	otherExtra.Title = "extra"
	otherExtra.Name = "wrong"
	if extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras whose name differ should not be equivalent")
	}

	otherExtra.Name = "extra"
	otherExtra.VerifiedSize = 9999
	if extra.IsEquivalentTo(&otherExtra, ChecksumValidationIfPresent, false) {
		t.Errorf("Extras whose verified size differ should not be equivalent")
	}
}

package manifest

import (
	"testing"
)

func TestManifestGameInstallerHasOneOfOses(t *testing.T) {
	installer := ManifestGameInstaller{
		Languages:     []string{"french"},
		Os:            "linux",
		Url:           "/dontknowdontcare",
		Title:         "installer",
		Name:          "installer",
		Version:       "vDontKnow",
		Date:          "2111-11-11",
		EstimatedSize: "1kb",
		VerifiedSize:  1000,
		Checksum:      "sdfdsfsdfdsfwe",
	}

	if installer.HasOneOfOses([]string{"windows", "macos"}) {
		t.Errorf("Should not indicate it has an os it doesn't have")
	}

	if !installer.HasOneOfOses([]string{"windows", "linux"}) {
		t.Errorf("Should indicate it has an os it has")
	}
}

func TestManifestGameInstallerHasOneOfLanguages(t *testing.T) {
	installer := ManifestGameInstaller{
		Languages:     []string{"french"},
		Os:            "linux",
		Url:           "/dontknowdontcare",
		Title:         "installer",
		Name:          "installer",
		Version:       "vDontKnow",
		Date:          "2111-11-11",
		EstimatedSize: "1kb",
		VerifiedSize:  1000,
		Checksum:      "sdfdsfsdfdsfwe",
	}

	if installer.HasOneOfLanguages([]string{"english", "german"}) {
		t.Errorf("Should not indicate it has a language it doesn't have")
	}

	if !installer.HasOneOfLanguages([]string{"english", "french"}) {
		t.Errorf("Should indicate it has a language it has")
	}
}

func TestManifestGameInstallerIsEquivalentTo(t *testing.T) {
	installer := ManifestGameInstaller{
		Languages:     []string{"french"},
		Os:            "linux",
		Url:           "/dontknowdontcare",
		Title:         "installer",
		Name:          "installer",
		Version:       "vDontKnow",
		Date:          "2111-11-11",
		EstimatedSize: "1kb",
		VerifiedSize:  1000,
		Checksum:      "icheck",
	}

	otherInstaller := ManifestGameInstaller{
		Languages:     []string{"english"},
		Os:            "windows",
		Url:           "/dontknowdontcare",
		Title:         "installer",
		Name:          "installer",
		Version:       "vDontKnowToo",
		Date:          "2111-12-12",
		EstimatedSize: "2kb",
		VerifiedSize:  1000,
		Checksum:      "icheck",
	}

	if !installer.IsEquivalentTo(&otherInstaller, false, false) {
		t.Errorf("Installers who url, title, name, verified size and checksum match should be equivalent when empty checksums are not tolerated")
	}

	if !installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers who url, title, name, verified size and checksum match should be equivalent when empty checksums are tolerated")
	}

	otherInstaller.Checksum = "idonotcheck"
	if installer.IsEquivalentTo(&otherInstaller, false, false) {
		t.Errorf("Installers whose checksum doesn't match should not be equivalent when empty checksums are not tolerated")
	}

	if installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers whose checksum doesn't match should not be equivalent when empty checksums are tolerated")
	}

	otherInstaller.Checksum = ""
	if installer.IsEquivalentTo(&otherInstaller, false, false) {
		t.Errorf("Installers whose checksum is empty should not be equivalent when empty checksums are not tolerated")
	}

	if !installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers whose checksum is empty should be equivalent when empty checksums are tolerated")
	}

	installer.Checksum = ""
	otherInstaller.Checksum = "icheck"
	if installer.IsEquivalentTo(&otherInstaller, false, false) {
		t.Errorf("Installers whose checksum is empty should not be equivalent when empty checksums are not tolerated")
	}

	if !installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers whose checksum is empty should be equivalent when empty checksums are tolerated")
	}

	installer.Checksum = ""
	otherInstaller.Checksum = ""
	if installer.IsEquivalentTo(&otherInstaller, false, false) {
		t.Errorf("Installers whose checksum is empty should not be equivalent when empty checksums are not tolerated")
	}

	if !installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers whose checksum is empty should be equivalent when empty checksums are tolerated")
	}

	otherInstaller.Title = "wrong"
	if installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers whose title differs should be not equivalent")
	}

	otherInstaller.Title = "installer"
	otherInstaller.Name = "wrong"
	if installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers whose name differs should be not equivalent")
	}

	otherInstaller.Name = "installer"
	otherInstaller.Url = "/dadadada"
	if installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers whose url differs should be not equivalent")
	}

	otherInstaller.Url = "/dontknowdontcare"
	otherInstaller.VerifiedSize = 9999
	if installer.IsEquivalentTo(&otherInstaller, true, false) {
		t.Errorf("Installers whose verified size differs should be not equivalent")
	}
}

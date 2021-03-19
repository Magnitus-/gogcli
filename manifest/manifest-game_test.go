package manifest

import (
	"testing"
)

func TestManifestGameImprintMissingChecksums(t *testing.T) {
	nextInstallers := []ManifestGameInstaller{
		ManifestGameInstaller{
			Languages: []string{"english"},
			Os: "windows",
			Url: "/dontknowdontcare",
			Title: "installer",
			Name: "installer",
			Version: "vDontKnowToo",
			Date: "2111-12-12",
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "",
		},
		ManifestGameInstaller{
			Languages: []string{"english"},
			Os: "windows",
			Url: "/dontknowdontcare2",
			Title: "installer2",
			Name: "installer2",
			Version: "vDontKnowToo",
			Date: "2111-12-12",
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "abcdefg",
		},
		ManifestGameInstaller{
			Languages: []string{"english"},
			Os: "windows",
			Url: "/dontknowdontcare3",
			Title: "installer3",
			Name: "installer3",
			Version: "vDontKnowToo",
			Date: "2111-12-12",
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "abcdefg",
		},
	}

	nextExtras := []ManifestGameExtra{
		ManifestGameExtra{
			Url: "/dontknowdontcare",
			Title: "extra",
			Name: "extra",
			Type: "ost",
			Info: 1,
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "",
		},
		ManifestGameExtra{
			Url: "/dontknowdontcare2",
			Title: "extra2",
			Name: "extra2",
			Type: "ost",
			Info: 1,
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "abcdefg",
		},
		ManifestGameExtra{
			Url: "/dontknowdontcare3",
			Title: "extra3",
			Name: "extra3",
			Type: "ost",
			Info: 1,
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "abcdefg",
		},
	}

	nextGame := ManifestGame{
		Id: 1,
		Title: "game",
		CdKey: "key",
		Tags: []string{"COMPLETED"},
		Installers: nextInstallers,
		Extras: nextExtras,
		EstimatedSize: "10mb",
		VerifiedSize: 10000,
	}

	prevInstallers := []ManifestGameInstaller{
		ManifestGameInstaller{
			Languages: []string{"english"},
			Os: "windows",
			Url: "/dontknowdontcare",
			Title: "installer",
			Name: "installer",
			Version: "vDontKnowToo",
			Date: "2111-12-12",
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "hijklmn", //Should imprint
		},
		ManifestGameInstaller{
			Languages: []string{"english"},
			Os: "windows",
			Url: "/dontknowdontcare2",
			Title: "installer2",
			Name: "installer2",
			Version: "vDontKnowToo",
			Date: "2111-12-12",
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "hijklmn", //Not empty in next, will not imprint
		},
		ManifestGameInstaller{
			Languages: []string{"english"},
			Os: "windows",
			Url: "/dontknowdontcare3",
			Title: "installer3",
			Name: "installer3",
			Version: "vDontKnowToo",
			Date: "2111-12-12",
			EstimatedSize: "1kb",
			VerifiedSize: 2000, //Won't match next, will not imprint
			Checksum: "hijklmn",
		},
	}

	prevExtras := []ManifestGameExtra{
		ManifestGameExtra{
			Url: "/dontknowdontcare",
			Title: "extra",
			Name: "extra",
			Type: "ost",
			Info: 1,
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "hijklmn", //Should imprint
		},
		ManifestGameExtra{
			Url: "/dontknowdontcare2",
			Title: "extra2",
			Name: "extra2",
			Type: "ost",
			Info: 1,
			EstimatedSize: "1kb",
			VerifiedSize: 1000,
			Checksum: "hijklmn", //Not empty in next, will not imprint
		},
		ManifestGameExtra{
			Url: "/dontknowdontcare3",
			Title: "extra3",
			Name: "extra3",
			Type: "ost",
			Info: 1,
			EstimatedSize: "1kb",
			VerifiedSize: 2000, //Won't match next, will not imprint
			Checksum: "hijklmn",
		},
	}

	prevGame := ManifestGame{
		Id: 1,
		Title: "game",
		CdKey: "key",
		Tags: []string{"COMPLETED"},
		Installers: prevInstallers,
		Extras: prevExtras,
		EstimatedSize: "10mb",
		VerifiedSize: 10000,
	}

	err := nextGame.ImprintMissingChecksums(&prevGame)
	if err != nil {
		t.Errorf("Imprinting two compatible games should work")
	}

	if nextGame.Installers[0].Checksum != "hijklmn" || nextGame.Installers[1].Checksum != "abcdefg" || nextGame.Installers[2].Checksum != "abcdefg" {
		t.Errorf("One of the game installers checksum is not as expected after imprinting")
	}

	if nextGame.Extras[0].Checksum != "hijklmn" || nextGame.Extras[1].Checksum != "abcdefg" || nextGame.Extras[2].Checksum != "abcdefg" {
		t.Errorf("One of the game extras checksum is not as expected after imprinting")
	}

	prevGame.Id = 2
	err = nextGame.ImprintMissingChecksums(&prevGame)
	if err == nil {
		t.Errorf("Imprinting two incompatible games should report an error")
	}
}

func TestManifestGameGetInstallerNamed(t *testing.T) {
	//TODO
}

func TestManifestGameGetExtraNamed(t *testing.T) {
	//TODO
}

func TestManifestGameTrimInstallers(t *testing.T) {
	//TODO
}

func TestManifestGameTrimExtras(t *testing.T) {
	//TODO
}

func TestManifestGameHasTitleTerms(t *testing.T) {
	//TODO
}

func TestManifestGameHasOneOfTags(t *testing.T) {
	//TODO
}

func TestManifestGameIsEmpty(t *testing.T) {
	//TODO
}

func TestManifestGameComputeVerifiedSize(t *testing.T) {
	//TODO
}

func TestManifestGameFillMissingFileInfo(t *testing.T) {
	//TODO
}
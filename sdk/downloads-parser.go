package sdk

import (
	"encoding/json"
	"errors"
)

/*
The json format of downloadable files returned by GOG takes the following form:
[
    <string representing language>,
    <objects representing downloads grouped by os>
]

Golang doesn't like that as the language only deals with arrays whose elements are of homogeneous type.

Hence, I have to do a little bit of wrangling here to deal with the above where I introduce a
strictly for parsing structure that can be either a language string or downloads.

Using lower-level parsing implementation, I detect whether or not an array element is a string.

If its a string, I set the value as a language and if its an object, I parse with the structure
expected by the files grouping.

The final expected parsing structure for each language entry is:
downloadFilesByLangParser [
	isDownloadFilesByOsOrIsLangParser{isLang=true, lang=...},
	isDownloadFilesByOsOrIsLangParser{
		isLang=false,
		filesByOs=downloadFilesByOsParser{
			Windows=[]downloadFileParser,
			Mac=[]downloadFileParser,
			Linux=[]downloadFileParser
		}
	}
]
*/

//We want this in the end
type gameDetailsDownloads []gameDetailsDownloadsFile

type gameDetailsDownloadsFile struct {
	Language  string
	Os        string
	ManualUrl string
	Name      string
	Version   string
	Date      string
	Size      string
}

//Everything below this is for parsing
func (g *gameDetailsDownloads) UnmarshalJSON(b []byte) error {
	var ds []downloadFilesByLangParser
	err := json.Unmarshal(b, &ds)
	if err != nil {
		return err
	}

	for _, d := range ds {
		if !d.isValid() {
			return errors.New("Downloads: Failure to parse downloads listing")
		}
		(*g) = append((*g), d.flatten()...)
	}

	return nil
}

type downloadFilesByLangParser []isDownloadFilesByOsOrIsLangParser

func (d downloadFilesByLangParser) isValid() bool {
	return len(d) == 2 && d[0].isLang && (!d[1].isLang)
}

func (d downloadFilesByLangParser) flatten() []gameDetailsDownloadsFile {
	var g []gameDetailsDownloadsFile
	for _, fs := range d[1].filesByOs.Windows {
		g = append(g, gameDetailsDownloadsFile{
			Language:  d[0].lang,
			Os:        "windows",
			ManualUrl: fs.ManualUrl,
			Name:      fs.Name,
			Version:   fs.Version,
			Date:      fs.Date,
			Size:      fs.Size,
		})
	}
	for _, fs := range d[1].filesByOs.Mac {
		g = append(g, gameDetailsDownloadsFile{
			Language:  d[0].lang,
			Os:        "mac",
			ManualUrl: fs.ManualUrl,
			Name:      fs.Name,
			Version:   fs.Version,
			Date:      fs.Date,
			Size:      fs.Size,
		})
	}
	for _, fs := range d[1].filesByOs.Linux {
		g = append(g, gameDetailsDownloadsFile{
			Language:  d[0].lang,
			Os:        "linux",
			ManualUrl: fs.ManualUrl,
			Name:      fs.Name,
			Version:   fs.Version,
			Date:      fs.Date,
			Size:      fs.Size,
		})
	}
	return g
}

type isDownloadFilesByOsOrIsLangParser struct {
	isLang    bool
	lang      string
	filesByOs downloadFilesByOsParser
}

type downloadFilesByOsParser struct {
	Windows []downloadFileParser
	Mac     []downloadFileParser
	Linux   []downloadFileParser
}

type downloadFileParser struct {
	ManualUrl string
	Name      string
	Version   string
	Date      string
	Size      string
}

func (i *isDownloadFilesByOsOrIsLangParser) UnmarshalJSON(b []byte) error {
	if string(b[0]) == "\"" || string(b[0]) == "'" {
		i.isLang = true
		i.lang = string(b[1 : len(b)-1])
	} else {
		var d downloadFilesByOsParser
		err := json.Unmarshal(b, &d)
		if err != nil {
			return err
		}

		i.isLang = false
		i.filesByOs = d
	}
	return nil
}

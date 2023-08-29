package manifest

import (
	"regexp"
)

type ManifestFilter struct {
	Titles          []string
	Oses            []string
	Languages       []string
	Tags            []string
	Installers      bool
	Extras          bool
	ExtraTypes      []string
	SkipUrls        []string
	HasUrls         []string
	Intersections   []ManifestFilter
	hasUrlsRegexes  []*regexp.Regexp
	skipUrlsRegexes []*regexp.Regexp
}

func NewManifestFilter(titles []string, oses []string, languages []string, tags []string, installers bool, extras bool, extraTypes []string, skipUrls []string, hasUrls []string) ManifestFilter {
	newFilter := ManifestFilter{
		Titles:          titles,
		Oses:            oses,
		Languages:       languages,
		Tags:            tags,
		Installers:      installers,
		Extras:          extras,
		ExtraTypes:      extraTypes,
		SkipUrls:        skipUrls,
		HasUrls:         hasUrls,
		Intersections:   []ManifestFilter{},
		skipUrlsRegexes: []*regexp.Regexp{},
		hasUrlsRegexes:  []*regexp.Regexp{},
	}
	newFilter.compileSkipUrls()
	newFilter.compileHasUrls()
	return newFilter
}

func (f *ManifestFilter) AddSkipUrl(url string) {
	for _, skipUrl := range f.SkipUrls {
		if url == skipUrl {
			return
		}
	}
	f.SkipUrls = append(f.SkipUrls, url)
	f.compileSkipUrls()
}

type FilterSkipUrlFn func(string) bool

func (f *ManifestFilter) GetSkipUrlFn() FilterSkipUrlFn {
	fn := func(u string) bool {
		for _, skipRegex := range (*f).skipUrlsRegexes {
			if skipRegex.MatchString(u) {
				return true
			}
		}

		return false
	}

	return fn
}

type FilterHasUrlFn func(string) bool

func (f *ManifestFilter) GetHasUrlFn() FilterHasUrlFn {
	fn := func(u string) bool {
		for _, hasRegex := range (*f).hasUrlsRegexes {
			if hasRegex.MatchString(u) {
				return true
			}
		}

		return false
	}

	return fn
}

func (f *ManifestFilter) IsEmpty() bool {
	isEmpty := len((*f).Titles) == 0 && len((*f).Oses) == 0
	isEmpty = isEmpty && len((*f).Languages) == 0 && len((*f).Tags) == 0
	isEmpty = isEmpty && len((*f).ExtraTypes) == 0 && len((*f).Intersections) == 0
	isEmpty = isEmpty && len((*f).SkipUrls) == 0
	isEmpty = isEmpty && len((*f).HasUrls) == 0
	return isEmpty
}

func (f *ManifestFilter) Intersect(other ManifestFilter) {
	if len((*f).Titles) == 0 {
		(*f).Titles = other.Titles
		other.Titles = []string{}
	}
	if len((*f).Oses) == 0 {
		(*f).Oses = other.Oses
		other.Oses = []string{}
	}
	if len((*f).Languages) == 0 {
		(*f).Languages = other.Languages
		other.Languages = []string{}
	}
	if len((*f).Tags) == 0 {
		(*f).Tags = other.Tags
		other.Tags = []string{}
	}
	if len((*f).ExtraTypes) == 0 {
		(*f).ExtraTypes = other.ExtraTypes
		other.ExtraTypes = []string{}
	}
	if (*f).Installers {
		(*f).Installers = other.Installers
	}
	if (*f).Extras {
		(*f).Extras = other.Extras
	}
	if len((*f).SkipUrls) == 0 {
		(*f).SkipUrls = other.SkipUrls
		other.SkipUrls = []string{}
	}
	if len((*f).HasUrls) == 0 {
		(*f).HasUrls = other.HasUrls
		other.HasUrls = []string{}
	}
	if (!other.IsEmpty()) || len(other.Intersections) > 0 {
		intersections := other.Intersections
		other.Intersections = []ManifestFilter{}
		intersections = append(intersections, other)
		(*f).Intersections = intersections
	}
}

func (f *ManifestFilter) compileSkipUrls() {
	(*f).skipUrlsRegexes = []*regexp.Regexp{}
	for _, u := range (*f).SkipUrls {
		(*f).skipUrlsRegexes = append((*f).skipUrlsRegexes, regexp.MustCompile(u))
	}
}

func (f *ManifestFilter) compileHasUrls() {
	(*f).hasUrlsRegexes = []*regexp.Regexp{}
	for _, u := range (*f).HasUrls {
		(*f).hasUrlsRegexes = append((*f).hasUrlsRegexes, regexp.MustCompile(u))
	}
}
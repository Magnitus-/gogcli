package manifest

import (
	"regexp"
	"sync"
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
	hasUrlsOnce     sync.Once
	skipUrlsOnce    sync.Once
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
	return newFilter
}

func (f *ManifestFilter) Copy() *ManifestFilter {
	newFilter := NewManifestFilter(
		f.Titles,
		f.Oses,
		f.Languages,
		f.Tags,
		f.Installers,
		f.Extras,
		f.ExtraTypes,
		f.SkipUrls,
		f.HasUrls,
	)
	return &newFilter
}

func (f *ManifestFilter) AddSkipUrl(url string) *ManifestFilter {
	for _, skipUrl := range f.SkipUrls {
		if url == skipUrl {
			return f
		}
	}
	newFilter := f.Copy()
	newFilter.SkipUrls = append(newFilter.SkipUrls, url)
	return newFilter
}

type FilterSkipUrlFn func(string) bool

func (f *ManifestFilter) GetSkipUrlFn() FilterSkipUrlFn {
	(*f).skipUrlsOnce.Do(func(){
		f.compileSkipUrls()
	})

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
	(*f).hasUrlsOnce.Do(func(){
		f.compileHasUrls()
	})

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

func (f *ManifestFilter) Intersect(other ManifestFilter) *ManifestFilter {
	newFilter := f.Copy()
	otherCopy := other.Copy()

	if len(newFilter.Titles) == 0 {
		newFilter.Titles = otherCopy.Titles
		otherCopy.Titles = []string{}
	}
	if len(newFilter.Oses) == 0 {
		newFilter.Oses = otherCopy.Oses
		otherCopy.Oses = []string{}
	}
	if len(newFilter.Languages) == 0 {
		newFilter.Languages = otherCopy.Languages
		otherCopy.Languages = []string{}
	}
	if len(newFilter.Tags) == 0 {
		newFilter.Tags = otherCopy.Tags
		otherCopy.Tags = []string{}
	}
	if len(newFilter.ExtraTypes) == 0 {
		newFilter.ExtraTypes = otherCopy.ExtraTypes
		otherCopy.ExtraTypes = []string{}
	}
	if newFilter.Installers {
		newFilter.Installers = otherCopy.Installers
	}
	if newFilter.Extras {
		newFilter.Extras = otherCopy.Extras
	}
	if len(newFilter.SkipUrls) == 0 {
		newFilter.SkipUrls = otherCopy.SkipUrls
		otherCopy.SkipUrls = []string{}
	}
	if len(newFilter.HasUrls) == 0 {
		newFilter.HasUrls = otherCopy.HasUrls
		otherCopy.HasUrls = []string{}
	}
	if (!otherCopy.IsEmpty()) || len(otherCopy.Intersections) > 0 {
		intersections := otherCopy.Intersections
		otherCopy.Intersections = []ManifestFilter{}
		intersections = append(intersections, *otherCopy)
		newFilter.Intersections = intersections
	}

	return newFilter
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
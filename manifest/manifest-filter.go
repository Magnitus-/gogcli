package manifest

type ManifestFilter struct {
	Titles []string
	Oses []string
	Languages []string
	Tags []string
	Installers bool
	Extras bool
	ExtraTypes []string
	Intersections []ManifestFilter
}

func NewManifestFilter (titles []string, oses []string, languages []string, tags []string, installers bool, extras bool, extraTypes []string) ManifestFilter {
	return ManifestFilter{
		Titles: titles,
		Oses: oses,
		Languages: languages,
		Tags: tags,
		Installers: installers,
		Extras: extras,
		ExtraTypes: extraTypes,
		Intersections: []ManifestFilter{},
	}
}

func (f *ManifestFilter) IsEmpty() bool {
	isEmpty := len((*f).Titles) == 0 && len((*f).Oses) == 0
	isEmpty = isEmpty && len((*f).Languages) == 0 && len((*f).Tags) == 0
	isEmpty = isEmpty && len((*f).ExtraTypes) == 0 && len((*f).Intersections) == 0
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
	if (!other.IsEmpty()) || len(other.Intersections) > 0 {
		intersections := other.Intersections
		other.Intersections = []ManifestFilter{}
		intersections = append(intersections, other)
		(*f).Intersections = intersections
	}
}
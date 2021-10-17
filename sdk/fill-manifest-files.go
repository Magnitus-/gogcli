package sdk

import "gogcli/manifest"

func (s *Sdk) fillManifestFiles(m *manifest.Manifest, concurrency int, pause int, tolerateDangles bool) ([]error, []error) {
	installersMap := (*m).GetUrlMappedInstallers()
	installerUrls := make([]string, len(installersMap))
	idx := 0
	for k, _ := range installersMap {
		installerUrls[idx] = k
		idx++
	}

	downloadInfos, fileInfoErrs, danglingInstallerErrs := s.GetManyDownloadFileInfo(installerUrls, concurrency, pause, tolerateDangles)
	if len(fileInfoErrs) > 0 {
		return fileInfoErrs, danglingInstallerErrs
	}
	for _, downloadFileInfo := range downloadInfos {
		(*installersMap[downloadFileInfo.url]).Name = downloadFileInfo.name
		(*installersMap[downloadFileInfo.url]).Checksum = downloadFileInfo.checksum
		(*installersMap[downloadFileInfo.url]).VerifiedSize = downloadFileInfo.size
	}

	extrasMap := (*m).GetUrlMappedExtras()
	extraUrls := make([]string, len(extrasMap))
	idx = 0
	for k, _ := range extrasMap {
		extraUrls[idx] = k
		idx++
	}

	var danglingExtraErrs []error
	downloadInfos, fileInfoErrs, danglingExtraErrs = s.GetManyDownloadFileInfo(extraUrls, concurrency, pause, tolerateDangles)
	if len(fileInfoErrs) > 0 {
		return fileInfoErrs, append(danglingInstallerErrs, danglingExtraErrs...)
	}
	for _, downloadFileInfo := range downloadInfos {
		(*extrasMap[downloadFileInfo.url]).Name = downloadFileInfo.name
		(*extrasMap[downloadFileInfo.url]).Checksum = downloadFileInfo.checksum
		(*extrasMap[downloadFileInfo.url]).VerifiedSize = downloadFileInfo.size
	}

	return []error{}, append(danglingInstallerErrs, danglingExtraErrs...)
}

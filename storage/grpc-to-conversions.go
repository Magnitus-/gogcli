package storage

import (
	"gogcli/manifest"
    "gogcli/storagegrpc"
)

func ConvertOs(os string) storagegrpc.Os {
	switch os {
	case "linux":
		return storagegrpc.Os_LINUX
	case "windows":
		return storagegrpc.Os_WINDOWS
	case "mac":
		return storagegrpc.Os_MACOS
	default:
		return storagegrpc.Os_UNSPECIFIED
	}
}

func ConvertManifestFilter(filter manifest.ManifestFilter) *storagegrpc.ManifestFilter {
	conversion := storagegrpc.ManifestFilter{
		Titles: filter.Titles,
		Installers: filter.Installers,
		Extras: filter.Extras,
		ExtraTypes: filter.ExtraTypes,
		Oses: []storagegrpc.Os{},
		Intersections: []*storagegrpc.ManifestFilter{},
	}

	for _, os := range filter.Oses {
		conversion.Oses = append(conversion.Oses, ConvertOs(os))
	}

	for _, subfilter := range filter.Intersections {
		conversion.Intersections = append(conversion.Intersections, ConvertManifestFilter(subfilter))
	}

	return &conversion
}

func ConvertManifestOverview(man manifest.Manifest) *storagegrpc.ManifestOverview {
	conversion := storagegrpc.ManifestOverview{
		EstimatedSize: man.EstimatedSize,
		VerifiedSize: man.VerifiedSize,
		Filter: ConvertManifestFilter(man.Filter),
	}

	return &conversion
}

func ConvertManifestGameInstaller(installer manifest.ManifestGameInstaller) *storagegrpc.ManifestGameInstaller {
	conversion := storagegrpc.ManifestGameInstaller{
		Name: installer.Name,
		Title: installer.Title,
		Url: installer.Url,
		TargetOs: ConvertOs(installer.Os),
		Languages: installer.Languages,
		Version: installer.Version,
		Date: installer.Date,
		EstimatedSize: installer.EstimatedSize,
		VerifiedSize: installer.VerifiedSize,
		Checksum: installer.Checksum,
	}

	return &conversion
}

func ConvertManifestGameExtra(extra manifest.ManifestGameExtra) *storagegrpc.ManifestGameExtra {
	conversion := storagegrpc.ManifestGameExtra{
		Name: extra.Name,
		Title: extra.Title,
		Url: extra.Url,
		Type: extra.Type,
		Info: int64(extra.Info),
		EstimatedSize: extra.EstimatedSize,
		VerifiedSize: extra.VerifiedSize,
		Checksum: extra.Checksum,
	}

	return &conversion
}

func ConvertManifestGame(game manifest.ManifestGame) *storagegrpc.ManifestGame {
	conversion := storagegrpc.ManifestGame{
		Id: game.Id,
		Title: game.Title,
		CdKey: game.CdKey,
		Tags: game.Tags,
		Installers: []*storagegrpc.ManifestGameInstaller{},
		Extras: []*storagegrpc.ManifestGameExtra{},
		EstimatedSize: game.EstimatedSize,
		VerifiedSize: game.VerifiedSize,
	}

	for _, installer := range game.Installers {
		conversion.Installers = append(conversion.Installers, ConvertManifestGameInstaller(installer))
	}

	for _, extra := range game.Extras {
		conversion.Extras = append(conversion.Extras, ConvertManifestGameExtra(extra))
	}

	return &conversion
}

func ConvertFileAction(action manifest.FileAction) *storagegrpc.FileAction {
	conversion := storagegrpc.FileAction{
		Title: action.Title,
		Name: action.Name,
		Url: action.Url,
		Kind: action.Kind,
		Action: action.Action,
	}

	return &conversion
}

func ConvertGameAction(action manifest.GameAction) *storagegrpc.GameAction {
	conversion := storagegrpc.GameAction{
		Title: action.Title,
		Slug: action.Slug,
		Id: action.Id,
		Action: action.Action,
		InstallerActions: []*storagegrpc.FileAction{},
		ExtraActions: []*storagegrpc.FileAction{},
	}

	for _, fileAction := range action.InstallerActions {
		conversion.InstallerActions = append(conversion.InstallerActions, ConvertFileAction(fileAction))
	}

	for _, fileAction := range action.ExtraActions {
		conversion.ExtraActions = append(conversion.ExtraActions, ConvertFileAction(fileAction))
	}

	return &conversion
}

func ConvertGrpcConfigs(conf GrpcConfigs) *storagegrpc.GrpcConfigs {
	conversion := storagegrpc.GrpcConfigs{
		Endpoint: conf.Endpoint,
	}

	return &conversion
}

func ConvertS3Configs(conf S3Configs) *storagegrpc.S3Configs {
	conversion := storagegrpc.S3Configs{
		Endpoint: conf.Endpoint,
		Region: conf.Region,
		Bucket: conf.Bucket,
		Tls: conf.Tls,
		AccessKey: conf.AccessKey,
		SecretKey: conf.SecretKey,
	}

	return &conversion
}

func ConvertSource(src Source) *storagegrpc.Source {
	conversion := storagegrpc.Source{
		Type: src.Type,
		S3Params: ConvertS3Configs(src.S3Params),
		FsPath: src.FsPath,
		GrpcParams: ConvertGrpcConfigs(src.GrpcParams),
	}

	return &conversion
}

func ConvertGameInfo(info manifest.GameInfo) *storagegrpc.GameInfo {
	conversion := storagegrpc.GameInfo{
		Id: info.Id,
		Slug: info.Slug,
		Title: info.Title,
	}

	return &conversion
}

func ConvertFileInfo(info manifest.FileInfo) *storagegrpc.FileInfo {
	conversion := storagegrpc.FileInfo{
		Game: ConvertGameInfo(info.Game),
		Kind: info.Kind,
		Name: info.Name,
		Checksum: info.Checksum,
		Size: info.Size,
		Url: info.Url,
	}

	return &conversion
}

func ConvertFileInfoNoCheck(info manifest.FileInfo) *storagegrpc.FileInfoNoCheck {
	conversion := storagegrpc.FileInfoNoCheck{
		Game: ConvertGameInfo(info.Game),
		Kind: info.Kind,
		Name: info.Name,
		Url: info.Url,
	}

	return &conversion
}
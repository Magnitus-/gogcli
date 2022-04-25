package storage

import (
	"errors"
	"fmt"
	"gogcli/manifest"
    "gogcli/storagegrpc"

	"google.golang.org/grpc/status"
)

func ConvertGrpcOs(os storagegrpc.Os) string {
	switch os {
	case storagegrpc.Os_LINUX:
		return "linux"
	case storagegrpc.Os_WINDOWS:
		return "windows"
	case storagegrpc.Os_MACOS:
		return "mac"
	default:
		return "unknown"
	}
}

func ConvertGrpcError(err error) error {
	statusErr, ok := status.FromError(err)
	if ok {
		return errors.New(fmt.Sprintf("Grpc error (code %d): %s", statusErr.Code(), statusErr.Message()))
	}
	return err
}

func ConvertGrpcGameInfo(info *storagegrpc.GameInfo) manifest.GameInfo {
	return manifest.GameInfo{
		Id: info.GetId(),
		Slug: info.GetSlug(),
		Title: info.GetTitle(),
	}
}

func ConvertGrpcFileInfo(info *storagegrpc.FileInfo) manifest.FileInfo {
	return manifest.FileInfo{
		Game: ConvertGrpcGameInfo(info.GetGame()),
		Kind: info.GetKind(),
		Name: info.GetName(),
		Checksum: info.GetChecksum(),
		Size: info.GetSize(),
		Url: info.GetUrl(),
	}
}

func ConvertGrpcManifestFilter(filter *storagegrpc.ManifestFilter) manifest.ManifestFilter {
	conversion := manifest.ManifestFilter{
		Titles: filter.GetTitles(),
		Oses: []string{},
		Languages: filter.GetLanguages(),
		Tags: filter.GetTags(),
		Installers: filter.GetInstallers(),
		Extras: filter.GetExtras(),
		ExtraTypes: filter.GetExtraTypes(),
		Intersections: []manifest.ManifestFilter{},
	}

	for _, os := range filter.GetOses() {
		conversion.Oses = append(conversion.Oses, ConvertGrpcOs(os))
	}

	for _, subfilter := range filter.GetIntersections() {
		conversion.Intersections = append(conversion.Intersections, ConvertGrpcManifestFilter(subfilter))
	}

	return conversion
}

func ConvertGrpcManifestOverview(man *storagegrpc.ManifestOverview) manifest.Manifest {
	return manifest.Manifest{
		Games: []manifest.ManifestGame{},
		EstimatedSize: man.GetEstimatedSize(),
		VerifiedSize: man.GetVerifiedSize(),
		Filter: ConvertGrpcManifestFilter(man.GetFilter()),
	}
}

func ConvertGrpcManifestGameInstaller(installer *storagegrpc.ManifestGameInstaller) manifest.ManifestGameInstaller {
	return manifest.ManifestGameInstaller{
		Languages: installer.GetLanguages(),
		Os: ConvertGrpcOs(installer.GetTargetOs()),
		Url: installer.GetUrl(),
		Title: installer.GetTitle(),
		Name: installer.GetName(),
		Version: installer.GetVersion(),
		Date: installer.GetDate(),
		EstimatedSize: installer.GetEstimatedSize(),
		VerifiedSize: installer.GetVerifiedSize(),
		Checksum: installer.GetChecksum(),
	}
}

func ConvertGrpcManifestGameExtra(extra *storagegrpc.ManifestGameExtra) manifest.ManifestGameExtra {
	return manifest.ManifestGameExtra{
		Url: extra.GetUrl(),
		Title: extra.GetTitle(),
		Name: extra.GetName(),
		Type: extra.GetType(),
		Info: int(extra.GetInfo()),
		EstimatedSize: extra.GetEstimatedSize(),
		VerifiedSize: extra.GetVerifiedSize(),
		Checksum: extra.GetChecksum(),
	}
}

func ConvertGrpcManifestGame(game *storagegrpc.ManifestGame) manifest.ManifestGame {
	conversion := manifest.ManifestGame{
		Id: game.GetId(),
		Slug: game.GetSlug(),
		Title: game.GetTitle(),
		CdKey: game.GetCdKey(),
		Tags: game.GetTags(),
		Installers: []manifest.ManifestGameInstaller{},
		Extras: []manifest.ManifestGameExtra{},
		EstimatedSize: game.GetEstimatedSize(),
		VerifiedSize: game.GetVerifiedSize(),
	}

	for _, installer := range game.GetInstallers() {
		conversion.Installers = append(conversion.Installers, ConvertGrpcManifestGameInstaller(installer))
	}
	
	for _, extra := range game.GetExtras() {
		conversion.Extras = append(conversion.Extras, ConvertGrpcManifestGameExtra(extra))
	}

	return conversion
}

func ConvertGrpcFileAction(action *storagegrpc.FileAction) manifest.FileAction {
	return manifest.FileAction{
		Title: action.GetTitle(),
		Name: action.GetName(),
		Url: action.GetUrl(),
		Kind: action.GetKind(),
		Action: action.GetAction(),
	}
}

func ConvertGrpcGameAction(action *storagegrpc.GameAction) manifest.GameAction {
	conversion := manifest.GameAction{
		Title: action.GetTitle(),
		Slug: action.GetSlug(),
		Id: action.GetId(),
		Action: action.GetAction(),
		InstallerActions: map[string]manifest.FileAction{},
		ExtraActions: map[string]manifest.FileAction{},
	}

	for _, fAction := range action.GetInstallerActions() {
		conversion.InstallerActions[fAction.GetName()] = ConvertGrpcFileAction(fAction)
	}

	for _, fAction := range action.GetExtraActions() {
		conversion.ExtraActions[fAction.GetName()] = ConvertGrpcFileAction(fAction)
	}

	return conversion
}

func ConvertGrpcGrpcConfigs(conf *storagegrpc.GrpcConfigs) GrpcConfigs {
	return GrpcConfigs{
		Endpoint: conf.GetEndpoint(),
	}
}

func ConvertGrpcS3Configs(conf *storagegrpc.S3Configs) S3Configs {
	return S3Configs{
		Endpoint: conf.GetEndpoint(),
		Region: conf.GetRegion(),
		Bucket: conf.GetBucket(),
		Tls: conf.GetTls(),
		AccessKey: conf.GetAccessKey(),
		SecretKey: conf.GetSecretKey(),
	}
}

func ConvertGrpcSource(src *storagegrpc.Source) Source {
	return Source{
		Type: src.GetType(),
		S3Params: ConvertGrpcS3Configs(src.GetS3Params()),
		FsPath: src.GetFsPath(),
		GrpcParams: ConvertGrpcGrpcConfigs(src.GetGrpcParams()),
	}
}
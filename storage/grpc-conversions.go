package storage

import (
	"errors"
	"fmt"
	"gogcli/manifest"
    "gogcli/storagegrpc"

	"google.golang.org/grpc/status"
)

func convertGrpcError(err error) error {
	statusErr, ok := status.FromError(err)
	if ok {
		return errors.New(fmt.Sprintf("Grpc error (code %d): %s", statusErr.Code(), statusErr.Message()))
	}
	return err
}

func convertGrpcGameInfo(info *storagegrpc.GameInfo) manifest.GameInfo {
	return manifest.GameInfo{
		Id: info.GetId(),
		Slug: info.GetSlug(),
		Title: info.GetTitle(),
	}
}

func convertGrpcFileInfo(info *storagegrpc.FileInfo) manifest.FileInfo {
	return manifest.FileInfo{
		Game: convertGrpcGameInfo(info.GetGame()),
		Kind: info.GetKind(),
		Name: info.GetName(),
		Checksum: info.GetChecksum(),
		Size: info.GetSize(),
		Url: info.GetUrl(),
	}
}
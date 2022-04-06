package storage

import (
	"context"
	"gogcli/manifest"
    "gogcli/storagegrpc"
	"io"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/codes"
)

type GrpcStore struct {
	Endpoint  string
	Connection *grpc.ClientConn
	Client storagegrpc.StorageServiceClient
}

func getGrpcStore(endpoint string) (GrpcStore, error) {
	conn, err := grpc.Dial(endpoint)
	if err != nil {
		return GrpcStore{}, err
	}

	return GrpcStore{endpoint, conn, storagegrpc.NewStorageServiceClient(conn)}, nil
}

func (g GrpcStore) GetListing() (*StorageListing, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listing := NewEmptyStorageListing(GrpcStoreDownloader{g})

	req := &storagegrpc.GetListingRequest{}
	resStream, err := g.Client.GetListing(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return &listing, err
	}
	
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break;
		}

		if err != nil {
			err = convertGrpcError(err)
			return &listing, err
		}

		listingGame := msg.GetListingGame()
		
		game := convertGrpcGameInfo(listingGame.GetGame())

		installers := []manifest.FileInfo{}
		for _, installer := range listingGame.GetInstallers() {
			installers = append(installers, convertGrpcFileInfo(installer))
		}

		extras := []manifest.FileInfo{}
		for _, extra := range listingGame.GetExtras() {
			extras = append(extras, convertGrpcFileInfo(extra))
		}
		
		listing.Games[game.Id] = StorageListingGame{
			Game: game,
			Installers: installers,
			Extras: extras,
		}
	}

	return &listing, nil
}

func (g GrpcStore) SupportsReaderAt() bool {
	return true
}

func (g GrpcStore) IsSelfValidating() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.IsSelfValidatingRequest{}
	res, err := g.Client.IsSelfValidating(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return false, err
	}

	return res.GetIsSelfValidating(), nil
}

func (g GrpcStore) GenerateSource() *Source {
	return nil
}

func (g GrpcStore) GetPrintableSummary() string {
	return ""
}

func (g GrpcStore) Exists() (bool, error) {
	return false, nil
}

func (g GrpcStore) Initialize() error {
	return nil
}

func (g GrpcStore) HasManifest() (bool, error) {
	return false, nil
}

func (g GrpcStore) HasActions() (bool, error) {
	return false, nil
}

func (g GrpcStore) HasSource() (bool, error) {
	return false, nil
}

func (g GrpcStore) StoreManifest(m *manifest.Manifest) error {
	return nil
}

func (g GrpcStore) StoreActions(a *manifest.GameActions) error {
	return nil
}

func (g GrpcStore) StoreSource(s *Source) error {
	return nil
}

func (g GrpcStore) LoadManifest() (*manifest.Manifest, error) {
	return nil, nil
}

func (g GrpcStore) LoadActions() (*manifest.GameActions, error) {
	return nil, nil
}

func (g GrpcStore) LoadSource() (*Source, error) {
	return nil, nil
}

func (g GrpcStore) RemoveActions() error {
	return nil
}

func (g GrpcStore) RemoveSource() error {
	return nil
}

func (g GrpcStore) AddGame(game manifest.GameInfo) error {
	return nil
}

func (g GrpcStore) RemoveGame(game manifest.GameInfo) error {
	return nil
}

func (g GrpcStore) UploadFile(source io.ReadCloser, file manifest.FileInfo) (string, error) {
	return "", nil
}

func (g GrpcStore) RemoveFile(file manifest.FileInfo) error {
	return nil
}

func (g GrpcStore) DownloadFile(file manifest.FileInfo) (io.ReadCloser, int64, error) {
	return nil, 0, nil
}

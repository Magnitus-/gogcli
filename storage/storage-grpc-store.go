package storage

import (
	"context"
	"gogcli/manifest"
    "gogcli/storagegrpc"
	"io"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/codes"
)

type GrpcConfigs struct {
	Endpoint string
}

type GrpcStore struct {
	configs *GrpcConfigs
	connection *grpc.ClientConn
	client storagegrpc.StorageServiceClient
}

func getGrpcStore(endpoint string) (GrpcStore, error) {
	conn, err := grpc.Dial(endpoint)
	if err != nil {
		return GrpcStore{}, err
	}
	configs := GrpcConfigs{endpoint}

	return GrpcStore{&configs, conn, storagegrpc.NewStorageServiceClient(conn)}, nil
}

func (g GrpcStore) GetListing() (*StorageListing, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listing := NewEmptyStorageListing(GrpcStoreDownloader{g})

	req := &storagegrpc.GetListingRequest{}
	resStream, err := g.client.GetListing(ctx, req)
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
	res, err := g.client.IsSelfValidating(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return false, err
	}

	return res.GetIsSelfValidating(), nil
}

func (g GrpcStore) GenerateSource() *Source {
	src := Source{
		Type:     "grpc",
		GrpcParams: (*g.configs),
	}
	return &src
}

func (g GrpcStore) GetPrintableSummary() (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.GetPrintableSummaryRequest{}
	res, err := g.client.GetPrintableSummary(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return "", err
	}

	return res.GetSummary(), nil
}

func (g GrpcStore) Exists() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.ExistsRequest{}
	res, err := g.client.Exists(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return false, err
	}

	return res.GetExists(), nil
}

func (g GrpcStore) Initialize() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.InitializeRequest{}
	_, err := g.client.Initialize(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) HasManifest() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.HasManifestRequest{}
	res, err := g.client.HasManifest(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return false, err
	}

	return res.GetHasManifest(), nil
}

func (g GrpcStore) HasActions() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.HasActionsRequest{}
	res, err := g.client.HasActions(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return false, err
	}

	return res.GetHasActions(), nil
}

func (g GrpcStore) HasSource() (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.HasSourceRequest{}
	res, err := g.client.HasSource(ctx, req)
	if err != nil {
		err = convertGrpcError(err)
		return false, err
	}

	return res.GetHasSource(), nil
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

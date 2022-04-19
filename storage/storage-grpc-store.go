package storage

import (
	"context"
	"errors"
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
		err = ConvertGrpcError(err)
		return &listing, err
	}
	
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break;
		}

		if err != nil {
			err = ConvertGrpcError(err)
			return &listing, err
		}

		listingGame := msg.GetListingGame()
		
		game := ConvertGrpcGameInfo(listingGame.GetGame())

		installers := []manifest.FileInfo{}
		for _, installer := range listingGame.GetInstallers() {
			installers = append(installers, ConvertGrpcFileInfo(installer))
		}

		extras := []manifest.FileInfo{}
		for _, extra := range listingGame.GetExtras() {
			extras = append(extras, ConvertGrpcFileInfo(extra))
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
		err = ConvertGrpcError(err)
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
		err = ConvertGrpcError(err)
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
		err = ConvertGrpcError(err)
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
		err = ConvertGrpcError(err)
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
		err = ConvertGrpcError(err)
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
		err = ConvertGrpcError(err)
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
		err = ConvertGrpcError(err)
		return false, err
	}

	return res.GetHasSource(), nil
}

func (g GrpcStore) StoreManifest(m *manifest.Manifest) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := g.client.StoreManifest(ctx)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	req := &storagegrpc.StoreManifestRequest{
		Manifest: &storagegrpc.Manifest{
			Content: &storagegrpc.Manifest_Overview{
				Overview: ConvertManifestOverview(*m),
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	for _, game := range (*m).Games {
		req := &storagegrpc.StoreManifestRequest{
			Manifest: &storagegrpc.Manifest{
				Content: &storagegrpc.Manifest_Game{
					Game: ConvertManifestGame(game),
				},
			},
		}
		err = stream.Send(req)
		if err != nil {
			err = ConvertGrpcError(err)
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) StoreActions(a *manifest.GameActions) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := g.client.StoreActions(ctx)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	for _, gameAction := range (*a) {
		req := &storagegrpc.StoreActionsRequest{
			GameAction: ConvertGameAction(gameAction),
		}
		err = stream.Send(req)
		if err != nil {
			err = ConvertGrpcError(err)
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) StoreSource(s *Source) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.StoreSourceRequest{
		Source: ConvertSource(*s),
	}
	_, err := g.client.StoreSource(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) LoadManifest() (*manifest.Manifest, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.LoadManifestRequest{}
	stream, err := g.client.LoadManifest(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return nil, err
	}

	res, resErr := stream.Recv()
	if resErr != nil {
		resErr = ConvertGrpcError(resErr)
		return nil, resErr
	}

	overview := res.GetManifest().GetOverview()
	if overview == nil {
		return nil, errors.New("Failure to get manifest with grpc store. Storage did not respect the established protocol of sending manifest overview first.")
	}
	man := ConvertGrpcManifestOverview(overview)


	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break;
		}

		if err != nil {
			err = ConvertGrpcError(err)
			return nil, err
		}

		game := res.GetManifest().GetGame()
		if game == nil {
			return nil, errors.New("Failure to get manifest with grpc store. Storage did not respect the established protocol of sending only manifest games after first message.")
		}

		man.Games = append(man.Games, ConvertGrpcManifestGame(game))
	}

	return &man, nil
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

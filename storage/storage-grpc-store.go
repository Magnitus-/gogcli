package storage

import (
	"bytes"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	actions := manifest.GameActions{}
	
	req := &storagegrpc.LoadActionsRequest{}
	stream, err := g.client.LoadActions(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return nil, err
	}
	
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break;
		}

		if err != nil {
			err = ConvertGrpcError(err)
			return nil, err
		}

		action := res.GetGameAction()
		convertedAction := ConvertGrpcGameAction(action)
		actions[convertedAction.Id] = convertedAction
	}

	return &actions, nil
}

func (g GrpcStore) LoadSource() (*Source, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.LoadSourceRequest{}
	res, err := g.client.LoadSource(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return nil, err
	}

	src := ConvertGrpcSource(res.GetSource())
	return &src, nil
}

func (g GrpcStore) RemoveActions() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.RemoveActionsRequest{}
	_, err := g.client.RemoveActions(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) RemoveSource() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.RemoveSourceRequest{}
	_, err := g.client.RemoveSource(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) AddGame(game manifest.GameInfo) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.AddGameRequest{
		Game: ConvertGameInfo(game),
	}
	_, err := g.client.AddGame(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) RemoveGame(game manifest.GameInfo) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.RemoveGameRequest{
		Game: ConvertGameInfo(game),
	}
	_, err := g.client.RemoveGame(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) RemoveFile(file manifest.FileInfo) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &storagegrpc.RemoveFileRequest{
		File: ConvertFileInfoNoCheck(file),
	}
	_, err := g.client.RemoveFile(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		return err
	}

	return nil
}

func (g GrpcStore) UploadFile(source io.ReadCloser, file manifest.FileInfo) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := g.client.UploadFile(ctx)
	if err != nil {
		err = ConvertGrpcError(err)
		return "", err
	}

	req := &storagegrpc.UploadFileRequest{
		Upload: &storagegrpc.FileUpload{
			Content: &storagegrpc.FileUpload_File{
				File: ConvertFileInfo(file),
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		err = ConvertGrpcError(err)
		return "", err
	}

	bufferSize := 1024
	readBuffer := make([]byte, bufferSize)
	var sendBuffer []byte
	var rLen int
	for err == nil {
		rLen, err = source.Read(readBuffer)
		if rLen > 0 {
			if rLen != bufferSize {
				sendBuffer = readBuffer[0:rLen-1]
			} else {
				sendBuffer = readBuffer
			}

			req := &storagegrpc.UploadFileRequest{
				Upload: &storagegrpc.FileUpload{
					Content: &storagegrpc.FileUpload_Data{
						Data: sendBuffer,
					},
				},
			}

			err = stream.Send(req)
			if err != nil {
				err = ConvertGrpcError(err)
				return "", err
			}
		}
	}

	if err != io.EOF {
		return "", err
	}

	res, closeErr := stream.CloseAndRecv()
	if closeErr != nil {
		closeErr = ConvertGrpcError(closeErr)
		return "", closeErr
	}

	return res.GetChecksum(), nil
}

type GrpcFileDownloader struct {
	ended bool
	accumulate *bytes.Buffer
	stream storagegrpc.StorageService_DownloadFileClient
	cancel context.CancelFunc
}

func (d *GrpcFileDownloader) Read(p []byte) (int, error) {
	if !(*d).ended {
		res, resErr := (*d).stream.Recv()
		if resErr != nil && resErr != io.EOF {
			resErr = ConvertGrpcError(resErr)
			return 0, resErr
		}

		if resErr == io.EOF {
			(*d).ended = true
		}

		data := res.GetDownload().GetData()
		if data == nil {
			return 0, errors.New("Failure to get data from grpc store. Storage did not respect the established protocol of sending data after first message.")
		}
	
		if len(data) > 0 {
			(*d).accumulate.Write(data)
		}
 	}

	var eofErr error
	if (*d).ended && (*d).accumulate.Len() <= len(p) {
		eofErr = io.EOF
	} else {
		eofErr = nil
	}

	readCount, readErr := (*d).accumulate.Read(p)
	if readErr != nil && readErr != io.EOF {
		return readCount, readErr
	}

	return readCount, eofErr
}

func (d *GrpcFileDownloader) Close() error {
	(*d).cancel()
	return nil
}

func (g GrpcStore) DownloadFile(file manifest.FileInfo) (io.ReadCloser, int64, error) {
	ctx, cancel := context.WithCancel(context.Background())

	req := &storagegrpc.DownloadFileRequest{
		File: ConvertFileInfo(file),
	}
	stream, err := g.client.DownloadFile(ctx, req)
	if err != nil {
		err = ConvertGrpcError(err)
		cancel()
		return nil, 0, err
	}

	res, resErr := stream.Recv()
	if resErr != nil {
		resErr = ConvertGrpcError(resErr)
		cancel()
		return nil, 0, resErr
	}

	expectedSize := res.GetDownload().GetExpectedSize()
	if expectedSize == 0 {
		cancel()
		return nil, 0, errors.New("Failure to get expected file size with grpc store. Storage did not respect the established protocol of sending expected size first.")
	}

	fileDownloader := GrpcFileDownloader{
		ended: false,
		accumulate: new(bytes.Buffer),
		stream: stream,
		cancel: cancel,
	}

	return &fileDownloader, expectedSize, nil
}

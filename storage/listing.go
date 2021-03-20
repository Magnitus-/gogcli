package storage

import (
	"crypto/md5"
	"encoding/hex"
	"gogcli/manifest"
	"io"
)

type StorageListingGame struct {
	Id         int64
	Installers []string
	Extras     []string
}

func NewEmptyStorageListingGame(id int64) StorageListingGame {
	return StorageListingGame{
		Id: id,
		Installers: make([]string, 0),
		Extras: make([]string, 0),
	}
}

type ListingFileRetrieval struct {
	File  manifest.FileInfo
	Error  error
}

func (g StorageListingGame) RetrieveFileInfo(name string, kind string, d Downloader, c chan ListingFileRetrieval) {
	handle, size, name, err := d.Download(
		g.Id,
		manifest.FileAction{
			Name: name,
			Kind: kind,
		},
	)
	defer handle.Close()
	if err != nil {
		c <- ListingFileRetrieval{
			File: manifest.FileInfo{},
			Error: err,
		}
		return
	}

	h := md5.New()
	io.Copy(h, handle)
	checksum := hex.EncodeToString(h.Sum(nil))
	c <- ListingFileRetrieval{
		File: manifest.FileInfo{
			GameId: g.Id,
			Kind: kind,
			Name: name,
			Checksum: checksum,
			Size: size,
		},
		Error: nil,
	}
}

type ListingGameRetrieval struct {
	Game  manifest.ManifestGame
	Error error
}

func (g StorageListingGame) RetrieveManifestGame(c chan ListingGameRetrieval, d Downloader) {
	var err error
	fileChan := make(chan ListingFileRetrieval)
	game := manifest.ManifestGame{
		Id: g.Id,
		Installers: make([]manifest.ManifestGameInstaller, 0),
		Extras: make([]manifest.ManifestGameExtra, 0),
	}

	for _, inst := range g.Installers {
		go g.RetrieveFileInfo(inst, "installer", d, fileChan)
	}
	for _, extr := range g.Extras {
		go g.RetrieveFileInfo(extr, "extra", d, fileChan)
	}

	filesCount := len(g.Installers) + len(g.Extras)
	for idx := 0; idx < filesCount; idx++ {
		fileRetrieval := <- fileChan
		if fileRetrieval.Error != nil {
			err = fileRetrieval.Error
		} else {
			fileInfo := fileRetrieval.File
			if fileInfo.Kind == "installer" {
				game.Installers = append(game.Installers, manifest.ManifestGameInstaller{
					Name: fileInfo.Name,
					VerifiedSize: fileInfo.Size,
					Checksum: fileInfo.Checksum,
				})
			} else {
				game.Extras = append(game.Extras, manifest.ManifestGameExtra{
					Name: fileInfo.Name,
					VerifiedSize: fileInfo.Size,
					Checksum: fileInfo.Checksum,				
				})
			}

		}
	}

	c <- ListingGameRetrieval{
		Game: game,
		Error: err,
	}
}

type StorageListing struct {
	Games map[int64]StorageListingGame
	downloads  Downloader
}

func (l *StorageListing) GetGameIds() []int64 {
	gameIds := make([]int64, len((*l).Games))
	idx := 0
	for id, _ := range (*l).Games {
		gameIds[idx] = id
		idx++
	}
	return gameIds
}

func NewEmptyStorageListing(d Downloader) StorageListing {
	return StorageListing{
		Games: map[int64]StorageListingGame{},
		downloads: d,
	}
}

func (l *StorageListing) GetManifest(concurrency int) (*manifest.Manifest, error) {
	m := manifest.NewEmptyManifest(manifest.ManifestFilter{})
	var err error
	gameChan := make(chan ListingGameRetrieval)
	processedGames := 0
	processingGames := 0
	gameIds := l.GetGameIds()
	for processedGames < len((*l).Games) {
		if err != nil && processingGames == 0 {
			return nil, err
		}
		someLeft := len(gameIds) > 0
		canLaunchMore := someLeft && processingGames < concurrency
		if canLaunchMore && err == nil {
			game := (*l).Games[gameIds[0]]
			gameIds = gameIds[1:]
			go game.RetrieveManifestGame(gameChan, l.downloads)
			processingGames++
		} else {
			gameRetrieval := <- gameChan
			if gameRetrieval.Error != nil {
				err = gameRetrieval.Error
			} else {
				(*m).Games = append((*m).Games, gameRetrieval.Game)
			}
			processedGames++
			processingGames--
		}
	}

	return m, err
}
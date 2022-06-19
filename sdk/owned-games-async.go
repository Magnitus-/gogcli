package sdk

import (
	"sync"
	"sync/atomic"
	"time"
)

type OwnedGamesPageReturn struct {
	page OwnedGamesPage
	err  error
}

func (s *Sdk) GetOwnedGamesPageAsync(page int, search string, returnVal chan OwnedGamesPageReturn) {
	o, err := s.GetOwnedGames(page, search)
	returnVal <- OwnedGamesPageReturn{page: o, err: err}
}

func (s *Sdk) GetAllOwnedGamesPagesSync(search string, concurrency int, pause int) ([]OwnedGamesPage, []error) {
	var pageCount int
	var currentPage int
	var pages []OwnedGamesPage
	var errs []error
	var callVal OwnedGamesPageReturn
	c := make(chan OwnedGamesPageReturn)

	go s.GetOwnedGamesPageAsync(1, search, c)
	callVal = <-c
	if callVal.err != nil {
		errs = append(errs, callVal.err)
		return pages, errs
	}

	if callVal.page.TotalPages == 0 {
		return pages, errs
	}

	pages = append(pages, callVal.page)
	pageCount = callVal.page.TotalPages
	currentPage = callVal.page.Page

	for currentPage < pageCount {
		maxPage := min(currentPage+concurrency, pageCount)

		for i := currentPage + 1; i <= maxPage; i++ {
			go s.GetOwnedGamesPageAsync(i, search, c)
		}

		for i := currentPage + 1; i <= maxPage; i++ {
			callVal = <-c
			if callVal.err != nil {
				errs = append(errs, callVal.err)
			} else {
				pages = append(pages, callVal.page)
			}
		}

		if len(errs) > 0 {
			return pages, errs
		}

		currentPage = maxPage
		time.Sleep(time.Duration(pause) * time.Millisecond)
	}
	return pages, errs
}

func (s *Sdk) GetAllOwnedGamesPages(done <-chan struct{}, search string, concurrency int, pause int) <-chan OwnedGamesPageReturn {
	var wg sync.WaitGroup
	ownedGamesPageChan := make(chan OwnedGamesPageReturn)

	page1, err := s.GetOwnedGames(1, search)
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case ownedGamesPageChan <- OwnedGamesPageReturn{page: page1, err: err}:
		case <-done:
			return
		}
	}()

	if err != nil || page1.TotalPages < 2 {
		go func() {
			wg.Wait()
			close(ownedGamesPageChan)
		}()
		return ownedGamesPageChan
	}

	totalPages := int64(page1.TotalPages)
	currentPage := int64(page1.Page)

	for idx := 0; idx < concurrency; idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for true {
				nextPage := atomic.AddInt64(&currentPage, 1)
				if nextPage > totalPages {
					break
				}
				page, err := s.GetOwnedGames(int(nextPage), search)
				select {
				case ownedGamesPageChan <- OwnedGamesPageReturn{page: page, err: err}:
				case <-done:
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ownedGamesPageChan)
	}()

	return ownedGamesPageChan
}
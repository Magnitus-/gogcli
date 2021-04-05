package sdk

import "time"

type ProductReturn struct {
	product Product
	dangling bool
	err  error
}

func (s *Sdk) GetProductAsync(gameId int64, returnVal chan ProductReturn) {
	p, d, err := s.GetProduct(gameId)
	if err != nil {
		p.Id = gameId
	}
	returnVal <- ProductReturn{product: p, dangling: d, err: err}
}

func (s *Sdk) GetManyProducts(gameIds []int64, concurrency int, pause int) ([]Product, []error, []error) {
	var errs []error
	var warnings []error
	var games []Product
	c := make(chan ProductReturn)

	i := 0
	for i < len(gameIds) {
		beginning := i
		target := min(len(gameIds), i+concurrency)
		for i < target {
			go s.GetProductAsync(gameIds[i], c)
			i++
		}

		y := beginning
		for y < target {
			returnVal := <-c
			if returnVal.err != nil {
				if returnVal.dangling {
					warnings = append(warnings, returnVal.err)
				} else {
					errs = append(errs, returnVal.err)
				}
			} else {
				games = append(games, returnVal.product)
			}
			y++
		}

		if len(errs) > 0 {
			return games, errs, warnings
		}

		if i < len(gameIds) {
			time.Sleep(time.Duration(pause) * time.Millisecond)
		}
	}

	return games, errs, warnings
}

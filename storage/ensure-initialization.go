package storage

func EnsureInitialization(s Storage) error {
	initialized, err := s.Exists()
	if err != nil {
		return err
	}

	if !initialized {
		err = s.Initialize()
		return err
	}

	return nil
}

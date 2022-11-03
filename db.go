package buckis

type db struct {
	*dict

	// maybe a place for indexes not sure yet
}

func DB() *db {
	d := newDict()

	go d.listenForCommands()

	go d.backgroundLoad()

	d.Load()

	return &db{
		d,
	}
}

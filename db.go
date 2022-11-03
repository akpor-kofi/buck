package buckis

type Db struct {
	*dict

	// maybe a place for indexes not sure yet
}

func DB() *Db {
	d := newDict()

	go d.listenForCommands()

	go d.backgroundLoad()

	d.Load()

	return &Db{
		d,
	}
}

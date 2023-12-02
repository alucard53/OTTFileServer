package files

import "log"

type File struct {
	route string
	id    string
}

type Files []File

var FileList = Files{
	{
		route: "./files/eot.mkv",
		id:    "a1",
	},
	{
		route: "./files/loki.mkv",
		id:    "a2",
	},
	{
		route: "./files/bm.mkv",
		id:    "a3",
	},
}

func (f Files) Search(id string) (string, error) {
	for _, v := range f {
		if v.id == id {
			log.Println(v.route)
			return v.route, nil
		}
	}

	return "", FileNotFound
}

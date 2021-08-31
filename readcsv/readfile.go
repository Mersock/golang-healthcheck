package readcsv

import (
	"encoding/csv"
	"io"
	"os"
)

type readerCSV struct {
	Filepath string
}

type ReadCSV interface {
	ReaderCSV() (links []string, err error)
}

func NewReadCSV(filepath string) ReadCSV {
	return &readerCSV{
		Filepath: filepath,
	}
}

func (r *readerCSV) ReaderCSV() (links []string, err error) {
	f, err := os.Open(r.Filepath)
	if err != nil {
		return links, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return links, err
		}
		for value := range record {
			if record[value] != "" {
				links = append(links, record[value])
			}
		}
	}

	return links, nil
}

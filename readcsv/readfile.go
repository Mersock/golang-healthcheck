package readcsv

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type readerCSV struct {
	Filepath string
}

type ReadCSV interface {
	ReaderCSV() (links []string)
}

func NewReadCSV(filepath string) ReadCSV {
	return &readerCSV{
		Filepath: filepath,
	}
}

func (r *readerCSV) ReaderCSV() (links []string) {
	f, err := os.Open(r.Filepath)
	if err != nil {
		log.Fatal("Unable to read input file "+r.Filepath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Unable to parse file as CSV for "+r.Filepath, err)
		}
		for value := range record {
			links = append(links, record[value])
		}
	}

	return links
}

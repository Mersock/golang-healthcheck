package readcsv

import (
	"encoding/csv"
	"os"
	"testing"
)

func mockFile(filename string) string {
	data := []string{"test"}
	csvFile, _ := os.Create(filename)
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	_ = w.Write(data)
	defer w.Flush()
	return filename
}

func TestReaderCSVFound(t *testing.T) {
	filename := mockFile("test_reader.csv")
	r := NewReadCSV(filename)
	links, err := r.ReaderCSV()
	if err != nil {
		t.Errorf("Expected nil,Actual %v", err)
	}
	if len(links) == 0 {
		t.Errorf("Expected links greater than 0, Actual %v", len(links))
	}
	e := os.Remove(filename)
	if e != nil {
		t.Errorf("no such file %v", e)
	}
}

func TestReaderCSVNotFound(t *testing.T) {
	r := NewReadCSV("xx")
	link, err := r.ReaderCSV()
	if len(link) > 0 {
		t.Errorf("Expected link 0 len, Actual %v", len(link))
	}
	if err == nil {
		t.Errorf("Expected Error no such file, Actual %v", err)
	}
}

package main

import (
	"encoding/csv"
	"os"
	"sync"
)

type CsvFile struct {
	sync.Mutex
	wr   *csv.Writer
	file *os.File
}

func NewCsvFile() (*CsvFile, error) {
	// TODO: delete file or display error if files already exists
	f, err := os.Create("output.csv")
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(f)

	return &CsvFile{
		wr:   w,
		file: f,
	}, nil
}

func (c *CsvFile) WriteLine(l []string) error {
	c.Lock()
	if err := c.wr.Write(l); err != nil {
		return err
	}
	c.Unlock()

	return nil
}

func (c *CsvFile) Close() error {
	c.Lock()

	c.wr.Flush()
	c.file.Close()

	c.Unlock()

	if err := c.wr.Error(); err != nil {
		return err
	}

	return nil
}

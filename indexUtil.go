package main

import (
	"bytes"
	"io/ioutil"
	"sort"
	"strings"
)

// WordDocFreq is a map of document ID to its word frequency.
type WordDocFreq map[int]int

// Index is a map of word and WordDocFreq
type Index map[string]WordDocFreq

type Document struct {
	Text string
}

// Maximum char in a document
const maxParaLen int = 1000

// GenerateIndexFromGivenData divides the text into documents and creates weighted inverted index on words
func GenerateIndexFromGivenData(filename string) (Index, []string, []Document, error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, err
	}

	allLines := strings.Split(string(dat), "\n")

	var docs []Document
	idx := make(Index)

	line := 0
	counter := 0

	for line < len(allLines) {
		trimmedLine := strings.TrimSpace(allLines[line])
		var buffer bytes.Buffer
		for len(trimmedLine) > 0 {
			buffer.WriteString(trimmedLine)
			buffer.WriteString(" ")
			line++
			trimmedLine = strings.TrimSpace(allLines[line])
		}

		if len(buffer.String()) > 0 {
			newDocs := splitDocumentText(buffer.String())
			for _, doc := range newDocs {
				docs = append(docs, Document{Text: doc})
				idx.Add(counter, Document{Text: doc})
				counter++
			}
		}

		line++
	}

	return idx, getAllKeysFromIndex(idx), docs, nil
}

func splitDocumentText(text string) []string {
	var newDocText []string
	i := 0
	docLength := len(text)
	for i < docLength {
		start := i
		end := i + maxParaLen

		if end > docLength {
			end = docLength
		}

		for end < docLength && text[end] != '.' {
			end++
		}
		d := strings.TrimSpace(text[start:end])
		newDocText = append(newDocText, d)

		end++
		i = end
	}

	return newDocText
}

func (idx Index) Add(id int, doc Document) {
	for _, token := range analyze(doc.Text) {
		ids := idx[token]

		if ids != nil {
			ids[id]++
			continue
		}

		m := make(WordDocFreq)
		m[id] = 1
		idx[token] = m
	}
}

func getAllKeysFromIndex(idx Index) []string {
	var keys []string
	for k := range idx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

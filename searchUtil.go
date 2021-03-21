package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/agnivade/levenshtein"
)

type SearchResponse struct {
	FinalQuery []string   `json:"finalQuery"`
	Results    []string `json:"results"`
}

type void struct{}

var member void

// Converts the text to lower cased tokens
func analyze(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	return tokens
}

// Tokenizes the string and generates on alpanumeric characters
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// Converts the text to lower case
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

func HandleSearch(idx Index, indexKeys []string, docs []Document) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		query, ok := r.URL.Query()["q"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		finalResults := Search(idx, indexKeys, docs, query[0])
		bytes, err := json.Marshal(finalResults)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}

func Search(idx Index, indexKeys []string, docs []Document, query string) SearchResponse {
	query = strings.TrimSpace(query)
	if len(query) == 0 {
		return SearchResponse{[]string{}, []string{}}
	}

	var tokens []string
	if query[0] == '/' && query[len(query)-1] == '/' && len(query) > 2 { // Checks for regular expression
		tokens = findTokenUsingRegex(idx, indexKeys, strings.TrimSpace(query[1:len(query)-1]))
	} else {
		tokens = analyze(query)
	}

	var finalQuery []string
	r := make(WordDocFreq)

	// If the exact word is present in the doc, show that, else find the nearest similar word using levenshtein dist
	for _, token := range tokens {
		ids, ok := idx[token]
		if !ok {
			simWord := findSimilarWord(idx, indexKeys, token)
			// TODO: check if replace is working.
			// strings.ReplaceAll(finalQuery, token, simWord)
            token = simWord
			ids = idx[simWord]
		}

        finalQuery = append(finalQuery, token)
		for key, val := range ids {
			r[key] += val
		}
	}

	docIds := make([]int, 0, len(r))

	for k := range r {
		docIds = append(docIds, k)
	}

	sort.Slice(docIds, func(i, j int) bool {
		wi, wj := r[docIds[i]], r[docIds[j]]
		if wi == wj {
			return docIds[i] < docIds[j]
		}
		return wi > wj
	})

	result := make([]string, 0, len(docIds))
	for _, key := range docIds {
		result = append(result, docs[key].Text)
	}

	return SearchResponse{finalQuery, result}
}

func findSimilarWord(idx Index, indexKeys []string, query string) string {

	var levenshteinDist []int

	for _, k := range indexKeys {
		levenshteinDist = append(levenshteinDist, levenshtein.ComputeDistance(k, query))
	}

	min := levenshteinDist[0]
	minIndex := 0
	for i, v := range levenshteinDist {
		if v < min {
			min = v
			minIndex = i
		}
	}

	return indexKeys[minIndex]
}

func findTokenUsingRegex(idx Index, indexKeys []string, query string) []string {

	re, _ := regexp.Compile(query)
	var tokens []string

	for _, key := range indexKeys {
		if re.MatchString(key) {
			tokens = append(tokens, key)
		}
	}

	return tokens
}

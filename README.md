# ShakeSearch

The [current search](https://nidhidhamnani-search.herokuapp.com/) is designed by keeping the most commonly used search features in mind.

## Methodology

- Divided the text into documents for easy readability of results
- The documents are created using the following steps:

  - Split the text in paragraphs, paragraphs are identified using two consecutive newlines  
  - If the paragraph exceeds 1000 characters, again split the paragraph into sub-paragraphs of 1000 characters
  - While splitting the paragraph, if a sentence is left incomplete due to length contraint, then completed the sentence using the nearest full-stop  

- Along with the documents creation, simultaneously, created weighted inverted index on the documents by maintaining a map of documents in which the word occurs and it's frequency in the document
- At the end, the results are sorted by word frequency in the document (Note: For multi-word, the frequencies are added)

## Features Added

### Backend

- Case insensitive exact word search
- Multiword search: Multiple words can be searched at a given time and the results will be calculates using OR operator on the given words
- Fuzzy word search: If a word does not completely match the words present in the whole text, a nearest word is calculated based on levenstein's distance
- Regex search: A regular expression can also be given in the search using /<regex>/. Note, multiword search is disabled on regex

### Frontend

- In all cases, highlighting the matching strings for readability
- In case of fuzzy search, highlighting the word for which the results are calculated
- Displaying total number of documents/paragraphs the word was found in

## Future Scope (Sorted Desc by Priority)

- Improve the search performance by using a better method for indexing and using a cache for frequent queries
- Show only 15/20 results with pagination
- Display a summary of the paragraph rather than displaying the complete paragraph and an option to expand
- Perform stemming to show words similar or suggestions for the query
- Show multiple different suggestions for fuzzy word and allow user to choose the correct word
- Enable fuzzy search along with regex
- Allow custom AND and OR operator
- Use ML model to identity semantically similar words for a given query

## Search Engine

![Search Engine](./images/Search.png?raw=true "Search Engine")
![Search Engine](./images/SearchWithFuzzy.png?raw=true "Search Engine (Fuzzy Words)")
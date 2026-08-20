# Song Crawler

This repository is a data parsing tool utlizing the website [1001albumsgenerator.com](https://1001albumsgenerator.com). It does three main things: 
 
-  Attempts to fetch all available album urls from the main website. This is done via golang.org/x/net/html package. Parsing HTML nodes recursively allows us to extract the different 'href' elements. 
-  Generates a goroutine worker queue that fetches each url and parses the information into computable album data.
-  Writes the results to a .csv file. 

## Usage 

For now the program just has one line argument.
```bash
song-crawler -output='myfile.csv'
```
If specifying an output file inside another directory, this directory must exist before running the program.

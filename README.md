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

## Resulting Data Example
```csv
title,artist,year,spotify_url,youtube_url,apple_url,genres
Tragic Songs of Life,The Louvin Brothers,1956,spotify:album:5JZCCzamSeSAzsePPi632S,https://www.youtube.com/results?search_query=Tragic%20Songs%20of%20Life%20-%20The%20Louvin%20Brothers,https://music.apple.com/album/1615489964,Folk|Country
Elvis Presley,Elvis Presley,1956,spotify:album:7GXP5OhYyPVLmcVfO9Iqin,https://www.youtube.com/results?search_query=Elvis%20Presley%20-%20Elvis%20Presley,https://music.apple.com/album/671019373,Rock
```

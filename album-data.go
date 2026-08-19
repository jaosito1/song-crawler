package main

import (
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type AlbumData struct {
	Title  string
	Artist string
	Year   int

	SpotifyUrl *url.URL
	YoutubeUrl *url.URL
	AppleUrl   *url.URL

	Genres []string
}

func getAlbumData(n *html.Node) AlbumData {
	var album AlbumData

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			switch getAttribute(c, "class") {
			case "album-title":
				album.Title = getTextValue(c)

			case "album-artist":
				albumNode := findByTag(c, atom.A)
				album.Artist = getTextValue(albumNode)

			case "album-year":
				yearNode := findByTag(c, atom.A)
				val, err := strconv.ParseInt(getTextValue(yearNode), 10, 0)

				if err == nil {
					album.Year = int(val)
				}

			case "streaming-links":
				for v := range c.Descendants() {
					if isTag(v, atom.A) {
						raw := getAttribute(v, "href")

						u, err := url.Parse(raw)
						if err != nil {
							continue
						}

						if u.Hostname() == "" && strings.Contains(u.String(), "spotify") {
							album.SpotifyUrl = u
						}

						if u.Hostname() == "www.youtube.com" {
							album.YoutubeUrl = u
						}

						if u.Hostname() == "music.apple.com" {
							album.AppleUrl = u
						}
					}
				}

			case "album-metadata":
				var genres []string

				spanNode := findByTag(c, atom.Span)
				for v := range spanNode.Descendants() {
					if isTag(v, atom.A) {
						genres = append(genres, getTextValue(v))
					}
				}

				album.Genres = genres
			}
		}
	}

	return album
}

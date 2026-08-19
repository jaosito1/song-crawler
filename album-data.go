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
	Year   string

	SpotifyUrl *url.URL
	YoutubeUrl *url.URL
	AppleUrl   *url.URL

	Genres []string
}

func (d AlbumData) Csv() []string {
	// We'll put genres in the format of 'blues-rock|heavy-metal' to avoid writing multiple CSV lines
	genres := strings.Join(d.Genres, "|")

	row := []string{
		d.Title,
		d.Artist,
		d.Year,
		getUrlString(d.SpotifyUrl),
		getUrlString(d.YoutubeUrl),
		getUrlString(d.AppleUrl),
		genres,
	}

	return row
}

func getUrlString(u *url.URL) string {
	if u != nil {
		return u.String()
	}

	return ""
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
				year := getTextValue(yearNode)

				_, err := strconv.ParseInt(year, 10, 0)
				if err == nil {
					album.Year = year
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

package gosub

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Client interface for client functionality
type Client interface {
	Search() ([]Subtitle, error)
	Download()
}

// OpenSubtitle struct site
type OpenSubtitle struct {
}

// Search for subtitle files
func (opensub *OpenSubtitle) Search(sp OpenSubtitleSearchParameters) ([]Subtitle, error) {
	// send request
	resp, err := get(generateSearchURL(sp))
	defer resp.Body.Close()
	panicOnError(err)

	// check the response
	// todo:: check response header status
	body, err := ioutil.ReadAll(resp.Body)
	panicOnError(err)

	// if 200 convert to struct
	openSubResp := OpenSubtitleResponse{}
	err = json.Unmarshal(body, &openSubResp)
	panicOnError(err)

	// get what you want from the response
	subtitles := []Subtitle{}
	for _, opensubtitle := range openSubResp {
		subtitles = append(subtitles, Subtitle{
			DownloadLink: opensubtitle.SubDownloadLink,
			Language:     opensubtitle.LanguageName,
			SubFormat:    opensubtitle.SubFormat,
		})
	}
	return subtitles, nil
}

// Download subtitle
func (opensub *OpenSubtitle) Download(filepath, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func generateSearchURL(sp OpenSubtitleSearchParameters) string {
	url := "https://rest.opensubtitles.org/search"

	// MovieSize / Hash
	// Episode / season / imdbid
	// query -> just a text
	// sublangid
	url = fmt.Sprintf("%s/%s-%d", url, "moviebytesize", sp.moviebytesize)
	url = fmt.Sprintf("%s/%s-%s", url, "moviehash", sp.moviehash)
	url = fmt.Sprintf("%s/%s-%s", url, "sublanguageid", sp.sublanguageid)
	log.Println(url)
	return url
}

func get(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "TemporaryUserAgent")
	return client.Do(req)
}

// OpenSubtitleSearchParameters search parameter
type OpenSubtitleSearchParameters struct {
	episode       int
	imdbid        string //(always format it as sprintf("%07d", $imdb) - when using imdb you can add /tags-hdtv for example.
	moviebytesize int64
	moviehash     string //(should be always 16 character, must be together with moviebytesize)
	query         string //(use url_encode, make sure " " is converted to "%20")
	season        int
	sublanguageid string //(if ommited, all languages are returned)
	tag           string //(use url_encode, make sure " " is converted to "%20")
}

// OpenSubtitleResponse is the response from subtitle site
type OpenSubtitleResponse []struct {
	IDMovie          string      `json:"IDMovie"`
	IDMovieImdb      string      `json:"IDMovieImdb"`
	IDSubMovieFile   string      `json:"IDSubMovieFile"`
	IDSubtitle       string      `json:"IDSubtitle"`
	IDSubtitleFile   string      `json:"IDSubtitleFile"`
	ISO639           string      `json:"ISO639"`
	InfoFormat       string      `json:"InfoFormat"`
	InfoOther        string      `json:"InfoOther"`
	InfoReleaseGroup string      `json:"InfoReleaseGroup"`
	LanguageName     string      `json:"LanguageName"`
	MatchedBy        string      `json:"MatchedBy"`
	MovieByteSize    string      `json:"MovieByteSize"`
	MovieFPS         string      `json:"MovieFPS"`
	MovieHash        string      `json:"MovieHash"`
	MovieImdbRating  string      `json:"MovieImdbRating"`
	MovieKind        string      `json:"MovieKind"`
	MovieName        string      `json:"MovieName"`
	MovieNameEng     interface{} `json:"MovieNameEng"`
	MovieReleaseName string      `json:"MovieReleaseName"`
	MovieTimeMS      string      `json:"MovieTimeMS"`
	MovieYear        string      `json:"MovieYear"`
	QueryCached      int         `json:"QueryCached"`
	QueryNumber      string      `json:"QueryNumber"`
	QueryParameters  struct {
		Episode       int    `json:"episode"`
		Imdbid        string `json:"imdbid"`
		Season        int    `json:"season"`
		Sublanguageid string `json:"sublanguageid"`
	} `json:"QueryParameters"`
	Score               float64 `json:"Score"`
	SeriesEpisode       string  `json:"SeriesEpisode"`
	SeriesIMDBParent    string  `json:"SeriesIMDBParent"`
	SeriesSeason        string  `json:"SeriesSeason"`
	SubActualCD         string  `json:"SubActualCD"`
	SubAddDate          string  `json:"SubAddDate"`
	SubAuthorComment    string  `json:"SubAuthorComment"`
	SubAutoTranslation  string  `json:"SubAutoTranslation"`
	SubBad              string  `json:"SubBad"`
	SubComments         string  `json:"SubComments"`
	SubDownloadLink     string  `json:"SubDownloadLink"`
	SubDownloadsCnt     string  `json:"SubDownloadsCnt"`
	SubEncoding         string  `json:"SubEncoding"`
	SubFeatured         string  `json:"SubFeatured"`
	SubFileName         string  `json:"SubFileName"`
	SubForeignPartsOnly string  `json:"SubForeignPartsOnly"`
	SubFormat           string  `json:"SubFormat"`
	SubFromTrusted      string  `json:"SubFromTrusted"`
	SubHD               string  `json:"SubHD"`
	SubHash             string  `json:"SubHash"`
	SubHearingImpaired  string  `json:"SubHearingImpaired"`
	SubLanguageID       string  `json:"SubLanguageID"`
	SubLastTS           string  `json:"SubLastTS"`
	SubRating           string  `json:"SubRating"`
	SubSize             string  `json:"SubSize"`
	SubSumCD            string  `json:"SubSumCD"`
	SubSumVotes         string  `json:"SubSumVotes"`
	SubTSGroup          string  `json:"SubTSGroup"`
	SubTSGroupHash      string  `json:"SubTSGroupHash"`
	SubTranslator       string  `json:"SubTranslator"`
	SubtitlesLink       string  `json:"SubtitlesLink"`
	UserID              string  `json:"UserID"`
	UserNickName        string  `json:"UserNickName"`
	UserRank            string  `json:"UserRank"`
	ZipDownloadLink     string  `json:"ZipDownloadLink"`
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

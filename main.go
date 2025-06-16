package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/flytam/filenamify"
	"github.com/tkrajina/gpxgo/gpx"
)

func usage() {
	fmt.Println("Create a GPX files based on a certain data export format.")
	fmt.Println("Usage: go run . <folder name>")
	os.Exit(0)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type TrackDataRef struct {
	Type string `json:"__type"`
	Name string `json:"name"`
	Iso  string `json:"iso,omitempty"`
	Url  string `json:"url,omitempty"`
}

// track data structure for one track
type TrackData struct {
	Distance            float64      `json:"distance,omitempty"`
	AltitudeMax         float64      `json:"altitudeMax,omitempty"`
	SpeedAverage        float64      `json:"speedAverage,omitempty"`
	Language            string       `json:"language,omitempty"`
	Points              TrackDataRef `json:"points"`
	Duration            float64      `json:"duration,omitempty"`
	RatingRoadCondition int64        `json:"ratingRoadCondition,omitempty"`
	Altitudes           TrackDataRef `json:"altitudes"`
	TimeCreated         TrackDataRef `json:"timeCreated"`
	Speeds              TrackDataRef `json:"speeds"`
	Image               TrackDataRef `json:"image"`
	AltitudeIncline     float64      `json:"altitudeIncline,omitempty"`
	SpeedMax            float64      `json:"speedMax,omitempty"`
	Preview200          TrackDataRef `json:"preview200"`
	Timestamps          TrackDataRef `json:"dates"`
	RatingFun           int64        `json:"ratingFun,omitempty"`
	AltitudeDecline     float64      `json:"altitudeDecline,omitempty"`
	Tags                []string     `json:"tags,omitempty"`
	Preview400          TrackDataRef `json:"preview400"`
	RatingScenery       int64        `json:"ratingScenery,omitempty"`
	Name                string       `json:"name"`
	Comment             string       `json:"comment,omitempty"`
	AltitudeMin         float64      `json:"altitudeMin,omitempty"`
	TourFeed            bool         `json:"tourFeed,omitempty"`
}

// default empty author name
var Author string = ""

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// folder with export data, extracted
	dirName := os.Args[1]

	// read author name from folder name
	if len(strings.Split(dirName, " ")) > 2 {
		Author = strings.Split(dirName, " ")[2]
	}

	// read overview file out of export directory
	payload, err := os.ReadFile(dirName + "/track data.json")
	check(err)

	// all tracks availbale from the overview
	var Tracks []TrackData

	// read JSON all tracks, parse track structure with JSON file references
	err = json.Unmarshal(payload, &Tracks)
	check(err)

	// create a GPX file for each track
	for _, Track := range Tracks {
		createGpx(Track, dirName)
	}
}

// points JSON read structure
type Points struct {
	Points [][]float64 `json:"points"`
}

// elevations JSON read structure
type Elevations struct {
	Elevations []float64 `json:"altitudes"`
}

// timestamps JSON read structure
type Timestamps struct {
	Timestamps []int64 `json:"dates"`
}

func loadJson(dirName string, fileName string, uri string) error {
	if len(strings.Trim(uri, " ")) <= 0 {
		return nil
	}
	// check if file is already present, ship downloading in case
	info, err := os.Stat(dirName + "/" + fileName)
	if info != nil || os.IsExist(err) {
		return nil
	}

	// Make HTTP GET request
	resp, err := http.Get(uri)
	if err != nil {
		return fmt.Errorf("failed to fetch the file: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the local file
	out, err := os.Create(dirName + "/" + fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write the response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("JSON file saved as %s\n", fileName)
	return nil
}

// create one track GPX file
func createGpx(Track TrackData, dirName string) {
	var Gpx gpx.GPX

	// identify track with name, fill meta data
	Gpx.Name = Track.Name
	Gpx.Description = Track.Comment
	Gpx.AuthorName = Author
	startTime, err := time.Parse("2006-01-02T15:04:05.000Z", Track.TimeCreated.Iso)
	check(err)
	Gpx.Time = &startTime

	// create the output gpx file into same folder as source data, use track name for file
	// fix names to comply with file system semantics
	outputGPXfilename, err := filenamify.Filenamify(startTime.Format(time.RFC3339)+" "+Gpx.Name+".gpx", filenamify.Options{
		Replacement: "-",
	})
	check(err)

	gpxFile, err := os.Create(dirName + "/" + outputGPXfilename)
	check(err)
	fmt.Println(outputGPXfilename)
	defer gpxFile.Close()

	// read points from points JSON file
	err = loadJson(dirName, Track.Points.Name, Track.Points.Url)
	check(err)
	pointsFile, err := os.ReadFile(dirName + "/" + Track.Points.Name)
	check(err)
	var points Points
	err = json.Unmarshal(pointsFile, &points)
	check(err)
	numPoints := len(points.Points)

	// read elevations from elevations JSON file
	err = loadJson(dirName, Track.Altitudes.Name, Track.Altitudes.Url)
	check(err)
	elevationsFile, err := os.ReadFile(dirName + "/" + Track.Altitudes.Name)
	check(err)
	var elevations Elevations
	err = json.Unmarshal(elevationsFile, &elevations)
	check(err)
	numElevations := len(elevations.Elevations)

	// read timestamps from timestamps JSON file
	err = loadJson(dirName, Track.Timestamps.Name, Track.Timestamps.Url)
	check(err)
	timestampsFile, err := os.ReadFile(dirName + "/" + Track.Timestamps.Name)
	check(err)
	var timestamps Timestamps
	err = json.Unmarshal(timestampsFile, &timestamps)
	check(err)
	numTimestamps := len(timestamps.Timestamps)

	if numPoints == numElevations && numPoints == numTimestamps {
		// all array are of same length
	} else {
		fmt.Printf("WARNING: Calimoto data size of your track not consistent for %s counting %s trackpoints: % d, %s elevations: %d, %s timestamps: %d . The GPX may represent your trip less accuratly, the bigger the difference is.", outputGPXfilename, Track.Points.Name, numPoints, Track.Altitudes.Name, numElevations, Track.Timestamps.Name, numTimestamps)
	}

	// number of points from number read in points JSON file, dimensions for elevation and timestamps are to be the same
	gpxPoints := make([]gpx.GPXPoint, len(points.Points))

	// combine threee source into one GPX point data item
	for i := range gpxPoints {
		gpxPoints[i].Point.Latitude = points.Points[i][0]
		gpxPoints[i].Point.Longitude = points.Points[i][1]
		if i < numElevations {
			gpxPoints[i].Point.Elevation = *gpx.NewNullableFloat64(float64(elevations.Elevations[i]))
		}
		if i < numTimestamps {
			gpxPoints[i].Timestamp = startTime.Add(time.Millisecond * time.Duration(timestamps.Timestamps[i]))
		}
	}

	// create a Track segment and set all points into it
	gpxSegment := gpx.GPXTrackSegment{
		Points: gpxPoints,
	}

	// create a track and set the track segment into it
	gpxTrack := gpx.GPXTrack{
		Name:     Track.Name,
		Comment:  Track.Comment,
		Segments: []gpx.GPXTrackSegment{gpxSegment},
	}

	// set the track into the GPX
	Gpx.Tracks = []gpx.GPXTrack{gpxTrack}

	// generate XML for the GPX file
	payload, err := Gpx.ToXml(gpx.ToXmlParams{Indent: true})
	check(err)

	// write GPX file with XML content
	_, err = gpxFile.Write(payload)
	check(err)
}

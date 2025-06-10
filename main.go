package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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
	Iso  string `json:"iso"`
}

// track data structure for one track
type TrackData struct {
	Distance            int64        `json:"distance"`
	AltitudeMax         int64        `json:"altitudeMax"`
	SpeedAverage        float64      `json:"speedAverage"`
	Language            string       `json:"language"`
	Points              TrackDataRef `json:"points"`
	Duration            int64        `json:"duration"`
	RatingRoadCondition int64        `json:"ratingRoadCondition"`
	Altitudes           TrackDataRef `json:"altitudes"`
	TimeCreated         TrackDataRef `json:"timeCreated"`
	Speeds              TrackDataRef `json:"speeds"`
	Image               TrackDataRef `json:"image"`
	AltitudeIncline     int64        `json:"altitudeIncline"`
	SpeedMax            float64      `json:"speedMax"`
	Preview200          TrackDataRef `json:"preview200"`
	Timestamps          TrackDataRef `json:"dates"`
	RatingFun           int64        `json:"ratingFun"`
	AltitudeDecline     int64        `json:"altitudeDecline"`
	Tags                []string     `json:"tags"`
	Preview400          TrackDataRef `json:"preview400"`
	RatingScenery       int64        `json:"ratingScenery"`
	Name                string       `json:"name"`
	Comment             string       `json:"comment"`
	AltitudeMin         int64        `json:"altitudeMin"`
	TourFeed            bool         `json:"tourFeed"`
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
	Elevations []int64 `json:"altitudes"`
}

// timestamps JSON read structure
type Timestamps struct {
	Timestamps []int64 `json:"dates"`
}

// replacer for potentially unsafe characters in file names
var fsReplacer = strings.NewReplacer("/", "-", "\\", "-", ":", "-")

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
	outputGPXfilename := fsReplacer.Replace(startTime.Format(time.RFC3339) + " " + Gpx.Name + ".gpx")
	gpxFile, err := os.Create(dirName + "/" + outputGPXfilename)
	check(err)
	fmt.Println(outputGPXfilename)
	defer gpxFile.Close()

	// read points from points JSON file
	pointsFile, err := os.ReadFile(dirName + "/" + Track.Points.Name)
	check(err)
	var points Points
	err = json.Unmarshal(pointsFile, &points)
	check(err)

	// read elevations from elevations JSON file
	elevationsFile, err := os.ReadFile(dirName + "/" + Track.Altitudes.Name)
	check(err)
	var elevations Elevations
	err = json.Unmarshal(elevationsFile, &elevations)
	check(err)

	// read timestamps from timestamps JSON file
	timestampsFile, err := os.ReadFile(dirName + "/" + Track.Timestamps.Name)
	check(err)
	var timestamps Timestamps
	err = json.Unmarshal(timestampsFile, &timestamps)
	check(err)

	// number of points from number read in points JSON file, dimensions for elevation and timestamps are to be the same
	gpxPoints := make([]gpx.GPXPoint, len(points.Points))

	// combine threee source into one GPX point data item
	for i := range gpxPoints {
		gpxPoints[i].Point.Latitude = points.Points[i][0]
		gpxPoints[i].Point.Longitude = points.Points[i][1]
		gpxPoints[i].Point.Elevation = *gpx.NewNullableFloat64(float64(elevations.Elevations[i]))
		gpxPoints[i].Timestamp = startTime.Add(time.Millisecond * time.Duration(timestamps.Timestamps[i]))
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

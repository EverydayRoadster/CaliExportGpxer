package main

import (
	"os"
	"fmt"
	"strings"
	"time"
	"encoding/json"
	"github.com/tkrajina/gpxgo/gpx"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type TrackDataRef struct {
	Type 	string `json:"__type"`
	Name 	string `json:"name"`
	Iso		string `json:"iso"`
}

type TrackData struct {
	Distance 		int64 `json:"distance"`
    AltitudeMax 	int64 `json:"altitudeMax"`
    SpeedAverage 	float64 `json:"speedAverage"`
    Language 		string `json:"language"`
    Points 			TrackDataRef `json:"points"`
    Duration		int64 `json:"duration"`
    RatingRoadCondition int64 `json:"ratingRoadCondition"`
    Altitudes    	TrackDataRef `json:"altitudes"` 
    TimeCreated		TrackDataRef `json:"timeCreated"`
    Speeds 			TrackDataRef `json:"speeds"`
    Image 			TrackDataRef `json:"image"`
    AltitudeIncline int64  `json:"altitudeIncline"`
    SpeedMax    	float64 `json:"speedMax"`
    Preview200    	TrackDataRef `json:"preview200"`
    Timestamps    	TrackDataRef `json:"dates"`
    RatingFun    	int64 `json:"ratingFun"`
    AltitudeDecline int64 `json:"altitudeDecline"`
    Tags			[]string `json:"tags"`
    Preview400    	TrackDataRef `json:"preview400"`
    RatingScenery 	int64 `json:"ratingScenery"`
    Name 		    string `json:"name"`
    Comment			string `json:"comment"`
    AltitudeMin    	int64 `json:"altitudeMin"`
    TourFeed		bool `json:"tourFeed"`
}

var Author string = ""


func main(){

	dirName := "data export CoreForce 2023-09-14"
	if len(os.Args) > 1 {
		dirName = os.Args[1]
	} 

	if len(strings.Split(dirName, " ")) > 2 {
		Author = strings.Split(dirName, " ")[2]
	}

	payload,err := os.ReadFile(dirName + "/track data.json" )
	check(err)

	var Tracks []TrackData

	err = json.Unmarshal(payload, &Tracks)
	check(err)

	for _,Track := range Tracks {
		createGpx(Track, dirName)
	}
}

type Points struct {
	Points 		[][]float64 `json:"points"`
}
type Elevations struct {
	Elevations 	[]int64 `json:"altitudes"`
}
type Timestamps struct {
	Timestamps []int64 `json:"dates"`
}

func createGpx(Track TrackData, dirName string){
	var Gpx gpx.GPX
 
	fmt.Println(Track.Name)
	Gpx.Name = Track.Name
	Gpx.Description = Track.Comment
	Gpx.AuthorName = Author
	startTime,err := time.Parse("2006-01-02T15:04:05.000Z", Track.TimeCreated.Iso)
	check(err)
	fmt.Println(startTime.String())
	Gpx.Time = &startTime

	gpxFile,err := os.Create(dirName + "/" + Gpx.Name + ".gpx")
	check(err)
	defer gpxFile.Close()

	// read points
	pointsFile,err := os.ReadFile(dirName + "/" + Track.Points.Name)
	check(err)
	var points Points
	err = json.Unmarshal(pointsFile, &points)
	check(err)

	// read elevations
	elevationsFile,err := os.ReadFile(dirName + "/" + Track.Altitudes.Name)
	check(err)
	var elevations Elevations
	err = json.Unmarshal(elevationsFile, &elevations)
	check(err)

	// read timestamps
	timestampsFile,err := os.ReadFile(dirName + "/" + Track.Timestamps.Name)
	check(err)
	var timestamps Timestamps
	err = json.Unmarshal(timestampsFile, &timestamps)
	check(err)


	gpxPoints := make([]gpx.GPXPoint, len(points.Points))

	for i,_ := range gpxPoints {
		gpxPoints[i].Point.Latitude = points.Points[i][0]
		gpxPoints[i].Point.Longitude = points.Points[i][1]
		gpxPoints[i].Point.Elevation = *gpx.NewNullableFloat64(float64(elevations.Elevations[i]))
		gpxPoints[i].Timestamp = startTime.Add(time.Millisecond * time.Duration(timestamps.Timestamps[i]))
	}

	gpxSegment := gpx.GPXTrackSegment{
		Points: gpxPoints,
	}

	gpxTrack := gpx.GPXTrack{
		Name: Track.Name,
		Comment: Track.Comment,
		Segments: []gpx.GPXTrackSegment{gpxSegment},
	}

	Gpx.Tracks = []gpx.GPXTrack{gpxTrack}

	payload, err := Gpx.ToXml( gpx.ToXmlParams{ Indent: true} )
	check(err)

	_,err = gpxFile.Write(payload)
	check(err)
}
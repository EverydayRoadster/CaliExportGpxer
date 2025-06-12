# CaliExportGpxer
A Go progrom to combine a data export, retrieved from or provided by Calimoto, and to convert it into proper GPX files for each trip.

## Use case
As user of the Calomoto service, I would like to have proper and complete GPX files of my trips, so that I can use them in many other applications subsequently. The GPX files should not only contain the track points, but also altitude information and timestamps when the trip was taken.

### Reasoning

GPX files downloaded from Calimoto website will contain point data only. They lack the association of points with a timestamp and elevation data. This is often not sufficient for further use of such a GPX file.

Since recording of GPX points in Calimoto is not done in fixed intervals (the duration between points recorded is flexible, based on vehicle velocity at any given location), the GPX data points can not easily be attributed with a fixed time interval. Also, the geolocation accuracy of the points recorded is not very precise, witch makes it unfavorable for later addition of elevation data based on a geographical model, especially not in "interesting" terrains.

This renders the GPX files, as they are downloadable from Calimoto, useless for many past-trip applications, e.g. odometer information rendered in video, a map animation rendered by GPX Animator application or usage in Relive.cc service etc.

However, Calimoto as a matter of fact, does store all the relevant data. Calimoto by default just does not expose it to the user for his own further use. 

Such personal data may be requested from a service provider, according to German law Art. 15 DSVGO. When asked for the personal data stored, Calimoto provides the user with an export of his personal data from their databases. 

A new version of this tool will accept an excerpt of the Calimoto Website application as input and download additional data ad hoc. A description of the procedure is given here: https://www.youtube.com/watch?v=EGufpKe_8h0

Unlike the GPX downloaded from the Calimoto web page directly, the export file provides not only point data but also a number of other metrics, like timestamps and elevations for each point. However, the data is stored as a number of time series database excerpts in JSON format (basically as the format provided by mapbox service, which in turn uses aws time series databases for storage).

This Go program converts the various JSON data sets of the export into a number of GPX files, one for each trip recorded.

## Usage

Preconditions:

1. The "track data.json" source file has been retrieved from your Calimoto route planner page or the export file has been received from Calimoto and extracted to a local folder. 
Hont: When extracting the .zip provided by Calimoto support, the user's email address is used as password for the protected ZIP file. Assuming the Zip file was provided as "data export EverydayRoadster 2022-04-01.zip", this will create a directory "data export EverydayRoadster 2022-04-01" containing a number of .json, .png and .jpg files. 
Hint: The export will not contain the picture files a user may have uploaded with a trip, which makes it not a complete report according to Art. 15 DSVGO. However, here we are not interested in those anyway and the images would be available for download from the website. 
2. The Go language is installed on the local computer. See https://go.dev/doc/install on how to do that.

Example:
```
go run . "data export EverydayRoadster 2022-04-01"
```
This example above would, for each trip contained in the export, create a .gpx file with the same name as a trip. The GPX file for each trip would contain the point data, timestamp for each point as well as elevation of the point as recorded.

## Legal
Art. 15 DSVGO allows an individual, in a reasonable manner, to gather a copy of his personal data stored with a service provider.

The text used to request this data from Calimoto was like:

"Ich möchte bitte Auskunft gem. Art. 15 DSGVO zu meinen bei Calimoto gespeicherten Fahrten. 
Ein Download meiner persönlichen Daten ist nicht komplett umfänglich aller meiner Daten.
Ein Download wird in der App durch Hinweis auf ein Premium-Abo verhindert. 
Die Inanspruchnahme des Auskunftsrechts ist jedoch vom Gesetzgeber grundsätzlich kostenlos vorgeschrieben. 
Herzlichen Dank!
Liebe Grüsse
<User> "
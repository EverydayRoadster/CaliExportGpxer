# CaliExportGpxer
A Go progrom to combine a data export provided by Calimoto, and to convert it into proper GPX files for each trip.

## Use case
Get a complete GPX file for a user's trips stored at Calimoto.

GPX download from Calimoto website will contain point data only, lacks the association of points with a timestamp and elevation data. Since recording of GPX points in Calimoto is not done in fixed intervals (duration between points recorded is flexible based on current velocity), the GPX data points can not easily be attributed with a fixed time interval. Also the accuracy of the points recorded is not very precise, witch makes it unfavorable for later addition of height based on a geographic model.

This renders the GPX files, as they are downloadable from Calimoto, useless for many past-trip applications, like a map animation rendered by GPX Animator application or Relive.cc.

However, Calimoto in fact does store all relevant data. Such personal data may be requested from a service provider, according to German law Art. 15 DSVGO. When asked for the personal data stored, Calimoto provides the user with an export of his personal data from their database. 

Unlike the GPX downloaded from the web page, the export provides not only point data but also a number of other metrics, like timestamps and elevations for each point. However, the data is stored as a number of time series database excerpts in JSON format (basically as the format provided by mapbox service, which in turn uses aws time series databases for storage).

This Go program converts the various JSON data sets of the export into a number of GPX files, one for each trip recorded.

## Usage

Preconditions:

1. The export file has been received from Calimoto and extracted to a local folder. When extracting, the user's email address is used as password for the protected ZIP file. Assuming the Zip file was provided as "data export EverydayRoadster 2022-04-01.zip", this will create a directory "data export EverydayRoadster 2022-04-01" containing a number of .json, .png and .jpg files. 
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
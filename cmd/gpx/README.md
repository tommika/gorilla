gpx Utility
===========

Summarize a GPX file (leverages [geo/gpx](../../geo/gpx) package)
```
$  ./bin/gpx ./geo/gpx/test-data/BeaconToColdSpring.gpx
<-- output -->
{
  "NumPoints": 954,
  "StartPt": {
    "Lat": 41.493350844830275,
    "Lon": -73.96022438071668,
    "Ele": 40.993408203125,
    "Time": "2015-12-31T10:52:58Z"
  },
  "EndPt": {
    "Lat": 41.42669535242021,
    "Lon": -73.96546960808337,
    "Ele": 15.5184326171875,
    "Time": "2015-12-31T16:11:09Z"
  },
  "Extent": {
    "MinLat": 41.42669535242021,
    "MaxLat": 41.4934463147074,
    "MinLon": -73.97041861899197,
    "MaxLon": -73.94215350039303
  },
  "Duration": 19091,
  "MinElePt": {
    "Lat": 41.42682728357613,
    "Lon": -73.96564370021224,
    "Ele": 10.231201171875,
    "Time": "2015-12-31T16:09:03Z"
  },
  "MaxElePt": {
    "Lat": 41.48140596225858,
    "Lon": -73.94462255761027,
    "Ele": 494.2550048828125,
    "Time": "2015-12-31T12:26:20Z"
  },
  "EleGain": 973.334228515625,
  "Distance": 12347.400339673362,
  "GridDistance": 12369.67370356709,
  "Area": 0,
  "Loopiness": 0,
  "IsClosed": false
}
```

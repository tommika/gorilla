<!--
Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
-->

gpx
===

Software for working with GPS data


### GPX (GPS Exchange Format)

Read a GPX document
```go
doc, err := gpx.ReadGpxDocument("./test-data/FortHillLoop.gpx")
```

Summarize it
```go
summary := doc.Summerize()
```

```json
{
  "NumPoints": 4335,
  "StartPt": {
    "Lat": 41.372543,
    "Lon": -73.930267,
    "Ele": 161.810898,
    "Time": "2024-02-11T17:49:42Z"
  },
  "EndPt": {
    "Lat": 41.372486,
    "Lon": -73.930338,
    "Ele": 159.915894,
    "Time": "2024-02-11T19:01:56Z"
  },
  "Extent": {
    "MinLat": 41.370436,
    "MaxLat": 41.376441,
    "MinLon": -73.930357,
    "MaxLon": -73.921727
  },
  "Duration": 4334,
  "MinElePt": {
    "Lat": 41.372486,
    "Lon": -73.930338,
    "Ele": 159.915894,
    "Time": "2024-02-11T19:01:56Z"
  },
  "MaxElePt": {
    "Lat": 41.372735,
    "Lon": -73.922482,
    "Ele": 259.486664,
    "Time": "2024-02-11T18:13:45Z"
  },
  "EleGain": 173.13690100000002,
  "Distance": 3240.9683927952315,
  "GridDistance": 3512.3636559531824,
  "Area": 447020.15903245256,
  "Loopiness": 137.9279600585404,
  "IsClosed": true
}
```

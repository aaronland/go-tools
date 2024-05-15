# go-geojson-tools

Common tools for describing and working with GeoJSON files.

## Tools

### bbox2feature

```
./bin/bbox2feature -bbox '37.220487,-122.632141,37.982092,-121.555481' -latlon | json_pp
{
   "geometry" : {
      "coordinates" : [
         [
            [
               -122.632141,
               37.220487
            ],
            [
               -121.555481,
               37.220487
            ],
            [
               -121.555481,
               37.982092
            ],
            [
               -122.632141,
               37.982092
            ],
            [
               -122.632141,
               37.220487
            ]
         ]
      ],
      "type" : "Polygon"
   },
   "properties" : {
      "bbox" : "37.220487,-122.632141,37.982092,-121.555481"
   },
   "type" : "Feature"
}
```

### bounds

```
$> ./bin/bounds -latlon -featurecollection /usr/local/sfomuseum/terminals.geojson
37.610449,-122.393473,37.621206,-122.380808
```

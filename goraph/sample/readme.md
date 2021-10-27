## Sample

### Common setting

Download [Open Street Map .osm.pbf file](https://download.geofabrik.de/asia.html) or GeoJSON LineString file.

```
$ curl https://download.geofabrik.de/asia/japan/kanto-latest.osm.pbf -o kanto-latest.osm.pbf
```

### Shortest path

1. Set Mapbox token

#### index.html

```js
  mapboxgl.accessToken = '[Mapbox token]';
  let map = new mapboxgl.Map({
    container: 'map',
    style: 'mapbox://styles/mapbox/streets-v9',
    zoom: 14,
    center: [139.7639394, 35.6840311]
  });
```

2. Then, start server

```
$ go run server.go
```

3. Access ``http://localhost:8000``

![](https://gyazo.com/57bfc8650083f8e3e7fef46c983b1a35.gif)

### Voronoi diagram


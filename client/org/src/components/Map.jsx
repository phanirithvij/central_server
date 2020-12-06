import { Icon } from "leaflet";
import "leaflet-geosearch/dist/geosearch.css";
import "leaflet-geosearch/dist/geosearch.umd";
import icon from "leaflet/dist/images/marker-icon.png";
import shadow from "leaflet/dist/images/marker-shadow.png";
import "leaflet/dist/leaflet";
import "leaflet/dist/leaflet.css";
import React, { useImperativeHandle, useMemo, useRef, useState } from "react";
import {
  MapContainer,
  Marker,
  Popup,
  TileLayer,
  useMapEvents,
} from "react-leaflet";
import "./Map.css";
import Search, { Address } from "./Search";

/* const provider = new OpenStreetMapProvider();
const control = new GeoSearchControl({
  provider,
  style: "bar",
  marker: {
    alt: "marker",
    draggable: true,
    icon: new Icon({ iconUrl: icon, shadowUrl: shadow }),
  },
  keepResult: true,
  searchLabel: "Search for address",
  popupFormat: ({ result, query }) => {
    console.log(query, result);
    return `${query.query}\n ${result.x}, ${result.y}`;
  },
});
 */

const LocationMarker = React.forwardRef((props, ref) => {
  const [position, setPosition] = useState(null);
  const [label, setLabel] = useState("Current location");
  const updateMarkerLabel = (tupl) => {
    setLabel("Fetcing marker location...");
    Address(tupl)
      .then((address) => {
        setLabel(address[0].label);
      })
      .catch((err) => {
        setLabel(err.toString());
        console.error(err);
      });
  };
  const map = useMapEvents({
    click(e) {
      let tupl = [e.latlng.lat, e.latlng.lng];
      setPosition(tupl);
      updateMarkerLabel(tupl);
      // map.locate();
    },
    locationfound(e) {
      map.flyTo(e.latlng);
      setPosition([e.latlng.lat, e.latlng.lng]);
      console.log(e.latlng);
    },
  });
  const markerRef = useRef(null);
  const eventHandlers = useMemo(
    () => ({
      dragend() {
        const marker = markerRef.current;
        if (marker != null) {
          const latlng = marker.getLatLng();
          let tupl = [latlng.lat, latlng.lng];
          setPosition(tupl);
          updateMarkerLabel(tupl);
        }
      },
    }),
    []
  );

  //   https://stackoverflow.com/a/61547777/8608146
  useImperativeHandle(
    ref,
    () => ({
      setPositionLabel: (latlong, label) => {
        setPosition(latlong);
        setLabel(label);
      },
      getLatLng: () => position,
    }),
    [setPosition, setLabel, position]
  );

  return position === null ? null : (
    <Marker
      draggable={true}
      eventHandlers={eventHandlers}
      ref={markerRef}
      //   https://gis.stackexchange.com/a/324925/173743
      icon={
        new Icon({
          iconUrl: icon,
          shadowUrl: shadow,
          iconSize: [25, 41],
          iconAnchor: [12, 41],
        })
      }
      position={position}
    >
      <Popup>
        {label} {position !== null && `${position[0]}, ${position[1]}`}
      </Popup>
    </Marker>
  );
});

function SearchWrapper(props) {
  const selectCallback = ({ x, y, label }) => {
    if (props.map) {
      const map = props.map;
      if (props.setPositionLabel) props.setPositionLabel([y, x], label);
      map.flyTo([y, x]);
    }
  };

  return <Search selectCallback={selectCallback} latlong={props.latlng} />;
}

function Map() {
  const [map, setMap] = useState(null);
  //   default iiit hyderabad
  const [center, setCenter] = useState([17.44511053681717, 78.34944901691728]);
  const childRef = useRef(null);
  const setPositionLabel = (x, l) => {
    childRef.current.setPositionLabel(x, l);
    setCenter(x);
  };

  const eventHandlers = useMemo(
    () => ({
      zoomend(e) {
        console.log("Zoom End");
        console.log(e);
      },
      zoom(e) {
        console.log("Zoom");
        console.log(e);
      },
      zoomstart(e) {
        console.log("Zoom Start");
        console.log(e);
      },
      zoomlevelschange(e) {
        console.log("Zoom Level");
        console.log(e);
      },
    }),
    []
  );

  return (
    <div className="mapwrap">
      <SearchWrapper
        map={map}
        setPositionLabel={setPositionLabel}
        latlong={childRef.current?.getLatLng()}
      />
      <MapContainer
        whenCreated={(m) => setMap(m)}
        className="map"
        center={center}
        eventHandlers={eventHandlers}
        zoom={13}
        scrollWheelZoom={true}
        placeholder={
          <div>Map not visible for some reason, try enable javascript?</div>
        }
      >
        <TileLayer
          attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <LocationMarker ref={childRef} />
      </MapContainer>
    </div>
  );
}

export default Map;

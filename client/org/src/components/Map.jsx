import { Icon } from "leaflet";
import "leaflet-geosearch/dist/geosearch.css";
import "leaflet-geosearch/dist/geosearch.umd";
import icon from "leaflet/dist/images/marker-icon.png";
import shadow from "leaflet/dist/images/marker-shadow.png";
import "leaflet/dist/leaflet";
import "leaflet/dist/leaflet.css";
import React, { useImperativeHandle, useMemo, useRef, useState } from "react";
import SVG from "react-inlinesvg";
import {
  MapContainer,
  Marker,
  Popup,
  TileLayer,
  useMap,
  useMapEvents,
} from "react-leaflet";
import copy from "./drawing.svg";
import "./Map.css";
import Search, { Address } from "./Search";

const POSITION_CLASSES = {
  bottomleft: "leaflet-bottom leaflet-left",
  bottomright: "leaflet-bottom leaflet-right",
  topleft: "leaflet-top leaflet-left",
  topright: "leaflet-top leaflet-right",
};

// https://stackoverflow.com/questions/31924890/leaflet-js-custom-control-button-add-text-hover
// https://react-leaflet.js.org/docs/example-react-control
function CurrentLocationControl({ position }) {
  const map = useMap();
  const positionClass =
    (position && POSITION_CLASSES[position]) || POSITION_CLASSES.topleft;
  return (
    <div className={positionClass} style={{ top: "75px" }}>
      <div
        style={{ cursor: "pointer" }}
        className="leaflet-control leaflet-bar center-flex-control"
        onClick={() => {
          map.locate();
        }}
      >
        <div className="locBtn"></div>
      </div>
    </div>
  );
}

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
      // check if location button was clicked
      let { x, y } = e.containerPoint;
      // get map's bounds
      let re = e.target._container.getBoundingClientRect();
      // get absolute position in the dom
      x += re.x;
      y += re.y;
      const target = document.elementFromPoint(x, y);
      if (
        target.classList.contains("locBtn") ||
        target.classList.contains("center-flex-control")
      ) {
        // Don't set marker and label if location button was clicked
        // locationfound will handle it
        return;
      }
      let tupl = [e.latlng.lat, e.latlng.lng];
      setPosition(tupl);
      updateMarkerLabel(tupl);
    },
    locationfound(e) {
      map.flyTo(e.latlng);
      let tupl = [e.latlng.lat, e.latlng.lng];
      setPosition(tupl);
      updateMarkerLabel(tupl);
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

  // https://stackoverflow.com/a/61547777/8608146
  // exposing child methods to parent component
  useImperativeHandle(
    ref,
    () => ({
      setPositionLabel: (latlong, label) => {
        setPosition(latlong);
        setLabel(label);
      },
      getLatLng: () => position,
      getLabel: () => label,
    }),
    [setPosition, setLabel, position, label]
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
      <Popup className="popupcl">
        <div className="popup-item">
          <span>{label}</span>
          <div
            className="iconbtn"
            title="Use as Address"
            onClick={() => props.copyCallback("label", label)}
          >
            {/* https://stackoverflow.com/a/41756265/8608146 */}
            <SVG className="svgicon" src={copy}>
              <div>use</div>
            </SVG>
          </div>
        </div>
        {/* if position exists show it on pop up*/}
        {position !== null && (
          <div className="popup-item">
            <span>{`${position[0]}, ${position[1]}`}</span>
            <div
              className="iconbtn"
              title="Use as Location"
              onClick={() => props.copyCallback("latlong", position)}
            >
              <SVG className="svgicon" src={copy}>
                <div>use</div>
              </SVG>
            </div>
          </div>
        )}
      </Popup>
    </Marker>
  );
});

function SearchWrapper(props) {
  // can't use useMap so passing as a prop
  const selectCallback = ({ x, y, label }) => {
    if (props.map) {
      const map = props.map;
      if (props.setPositionLabel) props.setPositionLabel([y, x], label);
      map.flyTo([y, x]);
    }
  };

  return <Search selectCallback={selectCallback} latlong={props.latlng} />;
}

function Map(props) {
  const [map, setMap] = useState(null);
  // default location is iiit hyderabad
  const [center, setCenter] = useState([17.44511053681717, 78.34944901691728]);
  const childRef = useRef(null);
  // sets marker position and label
  const setPositionLabel = (x, l) => {
    childRef.current.setPositionLabel(x, l);
    setCenter(x);
  };

  // need to do this because useMap doesn't work outside <MapContainer/>
  const onMapInit = (m) => {
    setMap(m);
  };

  return (
    <div className="mapwrap">
      <SearchWrapper
        map={map}
        setPositionLabel={setPositionLabel}
        latlong={childRef.current?.getLatLng()}
      />
      <MapContainer
        whenCreated={onMapInit}
        className="map"
        center={center}
        zoom={13}
        scrollWheelZoom={true}
        placeholder={
          <div>Map not visible for some reason, try enabling javascript?</div>
        }
      >
        <TileLayer
          attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />

        <CurrentLocationControl />
        <LocationMarker copyCallback={props.copyCallback} ref={childRef} />
      </MapContainer>
    </div>
  );
}

export default Map;

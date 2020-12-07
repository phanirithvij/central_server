import { Icon } from "leaflet";
import "leaflet-geosearch/dist/geosearch.css";
import "leaflet-geosearch/dist/geosearch.umd";
import icon from "leaflet/dist/images/marker-icon.png";
import shadow from "leaflet/dist/images/marker-shadow.png";
import "leaflet/dist/leaflet";
import "leaflet/dist/leaflet.css";
import React, {
  useEffect,
  useImperativeHandle,
  useMemo,
  useRef,
  useState,
} from "react";
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

function usePrevious(value) {
  const ref = useRef();
  useEffect(() => {
    ref.current = value;
  });
  return ref.current;
}

// https://stackoverflow.com/a/19746771/8608146
function cmparr(a, b) {
  return !a && !b && a.length === b.length && a.every((v, i) => v === b[i]);
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

  // https://stackoverflow.com/a/53446665/8608146
  const previous = usePrevious(position);

  // https://github.com/PaulLeCam/react-leaflet/issues/317#issuecomment-739856989
  const openPopup = () => {
    // TODO: bug Clicking on the map closes the popup automatically
    // not related to any custom logic => leaflet is the culprit
    // so not checking if open for now as a workaround
    // if (document.querySelector(".popupcl")) return; // popup is already open
    document.querySelectorAll(".marker-x")?.[1].click();
  };

  useEffect(() => {
    if (position !== null && !cmparr(position, previous)) openPopup();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [position]);

  // https://stackoverflow.com/a/61547777/8608146
  // exposing child methods to parent component
  useImperativeHandle(
    ref,
    () => ({
      setPositionLabel: (latlng, label) => {
        setPosition(latlng);
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
          className: "marker-x",
        })
      }
      position={position}
    >
      <Popup closeOnEscapeKey={true} className="popupcl">
        <div className="popup-item">
          <span>{label}</span>
          <div
            className="iconbtn"
            title="Use as Address"
            onClick={() => props.copyCallback("address", label)}
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
              onClick={() => props.copyCallback("location", position)}
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

  return <Search selectCallback={selectCallback} />;
}

function Map(props) {
  const [map, setMap] = useState(null);
  // default location is iiit hyderabad
  const [center, setCenter] = useState([17.4454957, 78.34854697544472]);
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
      <SearchWrapper map={map} setPositionLabel={setPositionLabel} />
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

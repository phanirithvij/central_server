import { lazy, useState } from "react";
import { Link } from "react-router-dom";
import Org from "../../models/org";
import "./index.css";

const Map = lazy(() => import("../../components/Map"));

export default function Register() {
  const [mapVis, setmapVis] = useState(false);

  const [org] = useState(new Org());

  const updateOrg = (e) => {
    console.log(e.target.name, e.target.value);
    if (e.target.name.startsWith("email")) {
      // form of email-0, email-1 etc.
      // 0 being primary
      let idx = e.target.name.split("-")[1];
      org["email"](e.target.value, idx);
    } else if (e.target.name === "confirm") {
      // check if password == confirm
      console.log(e.target.value, org._password);
    } else if (e.target.name === "location") {
      org["location"](e.target.value.split(",").map(e => parseFloat(e.trim())));
    } else {
      // set name = value
      org[e.target.name](e.target.value);
    }
  };

  return (
    <div>
      <h2>Register</h2>
      <Link to="/login">Login</Link>
      {/* 
        - Private bool
        - Name string
        - Email[] - emails for hub/user communication
        - ID string serverAssigned
        - Alias string - Human friendly org slug serverRecommended
        - Description string - human readable description
        - LocationStr string - Manual location address
        - Location
        - Server (servers ??)
      */}
      <form
        onChange={updateOrg}
        onSubmit={(e) => {
          e.preventDefault();
          org.create();
        }}
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <input type="text" name="name" placeholder="Name" />
        <input type="text" name="alias" placeholder="Alias" />
        <input type="text" name="email-0" placeholder="Email Primary" />
        <input type="text" name="email-1" placeholder="Email 1" />
        <input type="text" name="description" placeholder="Description" />
        <button onClick={() => setmapVis(!mapVis)}>
          {!mapVis ? "Show" : "Hide"} Map
        </button>
        <input type="text" name="address" placeholder="Address" />
        <input type="text" name="location" placeholder="Location Lat, Long" />
        <input type="password" name="password" placeholder="Password" />
        <input type="password" name="confirm" placeholder="Confirm password" />
        <button type="submit">Register</button>
      </form>
      {mapVis && <Map />}
    </div>
  );
}

import { lazy, useState } from "react";
import Org from "../../../models/org";
import "./index.css";

const Map = lazy(() => import("../../../components/Map"));

export default function Settings() {
  const [mapVis, setmapVis] = useState(false);
  const [passValid, setPassValid] = useState(false);
  const [pass, setPass] = useState();
  const [conf, setConf] = useState();
  const [org] = useState(new Org());
  const validatePass = () => {
    setPass(org.$password);
    setConf(org._confirm);
    if (!org.$password.startsWith(org._confirm)) {
      // entered wrong thing
      setPassValid(false);
      // console.log(org._confirm, "doesn't match", org.$password);
      return;
    } else {
      if (
        org._confirm.length >= org.$password.length &&
        org._confirm !== org.$password
      ) {
        setPassValid(false);
        // console.log(org._confirm, "doesn't match", org.$password);
        return;
      }
    }
    setPassValid(true);
  };
  const updateOrg = (e) => {
    if (e.target.name.startsWith("email")) {
      // form of email-0, email-1 etc.
      // 0 being primary
      let idx = e.target.name.split("-")[1];
      org["email"]({ email: e.target.value, private: false }, idx);
    } else if (e.target.name === "location") {
      if (e.target.value.split(",").length !== 2) {
        // TODO show error message
        console.log(e.target.value, "is not a valid location");
        return;
      }
      org["location"](
        e.target.value.split(",").map((e) => parseFloat(e.trim()))
      );
    } else if (e.target.name === "password" || e.target.name === "confirm") {
      // set name = value
      org[e.target.name](e.target.value);
      validatePass();
    } else {
      // set name = value
      org[e.target.name](e.target.value);
    }
  };
  // triggered when clicked on the copy icons in the marker popup
  const useCallback = (type, value) => {
    org[type](value);
    const inp = document.querySelector(`input[name="${type}"]`);
    inp.value = value.toString();
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (
      passValid &&
      org.$password &&
      org._confirm &&
      org.$password === org._confirm
    ) {
    } else {
      console.error("Non password");
      return;
    }
    if (org.$location && org.$location.length === 2) {
    } else {
      console.error("Non locaceon");
      return;
    }
    org
      .update()
      .then((res) => res.json())
      .then((data) => console.log(data))
      .catch((err) => console.error(err));
  };
  return (
    <div>
      <form
        onChange={updateOrg}
        onSubmit={handleSubmit}
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <input type="text" name="name" placeholder="Name" />
        <input type="text" name="alias" placeholder="Alias" />
        {/* TODO list of emails */}
        {/* TODO private property to emails */}
        <input type="text" name="email-0" placeholder="Email Primary" />
        <input type="text" name="email-1" placeholder="Email 1" />
        <input type="text" name="description" placeholder="Description" />
        <input type="text" name="address" placeholder="Address" />
        <input type="text" name="location" placeholder="Location Lat, Long" />
        <input type="checkbox" name="private" />
        <input type="password" name="password" placeholder="Password" />
        <input type="password" name="confirm" placeholder="Confirm password" />
        <label htmlFor="confirm">
          {!passValid &&
            pass &&
            conf &&
            `Passwords don't match ${pass} , ${conf}`}
        </label>
        <button type="submit">Register</button>
      </form>
      <button onClick={() => setmapVis(!mapVis)}>
        {!mapVis ? "Show" : "Hide"} Map
      </button>
      {mapVis && <Map copyCallback={useCallback} />}
    </div>
  );
}

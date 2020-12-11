import { lazy, useEffect, useState } from "react";
import Org from "../../../models/org";
import "./index.css";

// Must be lazy for it is ~ 1MB, gziped 200 KB
const Map = lazy(() => import("../../../components/Map"));

export default function Settings() {
  const [mapVis, setmapVis] = useState(false);
  const [passValid, setPassValid] = useState(false);
  const [pass, setPass] = useState();
  const [conf, setConf] = useState();
  const [org] = useState(new Org());

  const [emails, setEmails] = useState();

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

  useEffect(() => {
    org
      .fetch()
      .then(async (x) => {
        switch (x.status) {
          case 401:
            // TODO not logged in redirect to login page
            break;

          case 404:
            // TODO Logged in but org doesn't exist
            // Show user that warning and a logout button
            break;

          case 200:
            // Logged in and got the org details
            // Fill up org
            let json = await x.json();
            Object.keys(json).forEach((key) => {
              // console.log(key);
              if (key === "id") return;
              if (key === "emails") {
                org["emails"](json[key]);
                // console.log(org._emailList());
                setEmails(org._emailList());
                return;
              }
              // console.log(org[key], key);
              org[key](json[key]);
              if (typeof json[key] === "boolean") {
                document.querySelector(`input[name="${key}"]`).checked =
                  json[key];
              } else
                document.querySelector(`input[name="${key}"]`).value =
                  json[key];
            });
            break;

          default:
            break;
        }
      })
      .catch((err) => console.error(err));
  }, [org]);

  const updateOrg = (e) => {
    if (e.target.name.startsWith("email")) {
      // form of email-0, email-1 etc.
      // 0 being primary
      let parts = e.target.name.split("-");
      // TODO update email private
      if (parts.length === 2) {
        let idx = parts[1];
        // email
        org["email"]({ email: e.target.value }, idx);
      } else if (parts.length === 3) {
        let idx = parts[2];
        // private
        org["email"]({ private: e.target.checked }, idx);
      }
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
      org[e.target.name](
        e.target.type === "checkbox" ? e.target.checked : e.target.value
      );
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
    if (org.$password && org._confirm) {
      // we modfied the password
      if (passValid && org.$password === org._confirm) {
      } else {
        console.error("Non password");
        return;
      }
    }
    if (org.$location && org.$location.length === 2) {
    } else {
      console.error("Bad locaceon");
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
        {emails !== undefined &&
          emails.map((email, index) => (
            <div key={index}>
              {email.main && (
                <label htmlFor={`email-${index}`}>Main Email</label>
              )}
              <input
                type="text"
                name={`email-${index}`}
                placeholder={
                  index === 0 ? "Email Primary" : `Email ${index + 1}`
                }
                defaultValue={email.email}
              />
              <label htmlFor={`email-${index}`}>Private</label>
              <input
                type="checkbox"
                defaultChecked={email.private}
                name={`email-private-${index}`}
              />
            </div>
          ))}
        <input type="text" name="description" placeholder="Description" />
        <input type="text" name="address" placeholder="Address" />
        <input type="text" name="location" placeholder="Location Lat, Long" />
        <input type="checkbox" name="privateLoc" />
        <input type="checkbox" name="private" />
        <label htmlFor="oldPassword">Change Password</label>
        <input type="password" name="oldPassword" placeholder="Old Password" />
        <input type="password" name="newPassword" placeholder="Password" />
        <input type="password" name="confirm" placeholder="Confirm password" />
        <label htmlFor="confirm">
          {!passValid &&
            pass &&
            conf &&
            `Passwords don't match ${pass} , ${conf}`}
        </label>
        <button type="submit">Update</button>
      </form>
      <button onClick={() => setmapVis(!mapVis)}>
        {!mapVis ? "Show" : "Hide"} Map
      </button>
      {mapVis && <Map copyCallback={useCallback} />}
    </div>
  );
}

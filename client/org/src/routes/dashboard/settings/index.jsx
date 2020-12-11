import { Puff, useLoading } from "@agney/react-loading";
import Tippy from "@tippyjs/react";
import { lazy, useEffect, useState } from "react";
import SVG from "react-inlinesvg";
import "tippy.js/dist/tippy.css"; // optional
import Org from "../../../models/org";
import "./index.css";
import svgopen from "./map.svg";
import svgclose from "./mapclose.svg";
import AlertDismissible from "../../../components/Alert";
import { Button } from "antd";

// Must be lazy for it is ~ 1MB, gziped 200 KB
const Map = lazy(() => import("../../../components/Map"));

export default function Settings() {
  const [mapVis, setmapVis] = useState(false);
  const [passValid, setPassValid] = useState(false);
  const [pass, setPass] = useState();
  const [conf, setConf] = useState();
  const [org] = useState(new Org());
  const [clientValidError, setClientValidError] = useState();
  const [serverValidError, setServerValidError] = useState();
  const [sending, setSending] = useState();
  const [loggedin, setLoggedin] = useState();
  const [loggedinJson, setLoggedinJson] = useState();
  const [done, setDone] = useState();

  const [emails, setEmails] = useState();

  const { containerProps, indicatorEl } = useLoading({
    loading: true,
    indicator: <Puff width="50" />,
  });

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
    const json = loggedinJson;
    if (json === undefined) return;
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
        document.querySelector(`input[name="${key}"]`).checked = json[key];
      } else {
        if (["address", "description"].includes(key)) {
          document.querySelector(`textarea[name="${key}"]`).value = json[key];
        } else document.querySelector(`input[name="${key}"]`).value = json[key];
      }
    });
  }, [loggedinJson, org]);

  useEffect(() => {
    org
      .fetch()
      .then(async (x) => {
        switch (x.status) {
          case 401:
            // TODO not logged in redirect to login page
            setServerValidError("Not logged in, please login");
            break;

          case 404:
            // TODO Logged in but org doesn't exist
            // Show user that warning and a logout button
            setServerValidError(
              "Organization doesn't exist on our server please mail us if you believe this to be a mistake."
            );
            break;

          case 200:
            // Logged in and got the org details
            // Fill up org
            let json = await x.json();
            setLoggedin(true);
            setLoggedinJson(json);
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
    let elem = "input";
    if (type === "address") {
      elem = "textarea";
    }
    const inp = document.querySelector(`${elem}[name="${type}"]`);
    inp.value = value.toString();
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (org.$password && org._confirm) {
      // we modfied the password
      if (passValid && org.$password === org._confirm) {
      } else {
        setClientValidError("Password is not a valid password");
        return;
      }
    }
    if (org.$location && org.$location.length === 2) {
    } else {
      setClientValidError("Location is required");
      return;
    }
    setSending(true);
    org
      .update()
      .then(async (res) => {
        setSending(false);
        const jsonD = await res.json();
        switch (res.status) {
          case 422:
            console.error(jsonD["error"]);
            setServerValidError(jsonD["messages"].join("\n"));
            break;
          case 200:
            // successfully updted
            setDone(true);
            break;
          case 500:
            setServerValidError(jsonD["error"]);
            break;
          default:
            break;
        }
      })
      .catch((err) => {
        setSending(false);
        console.error(err);
      });
  };
  return (
    <div>
      {/* Org alias will be undefined only if not loggedin  */}
      {loggedin !== undefined && loggedin && (
        <div className="form-wrap">
          <form onChange={updateOrg} onSubmit={handleSubmit} className="formx">
            {org.$name !== undefined && org.$name === "" && (
              <AlertDismissible
                show
                content={"Please add all the details"}
                variant="info"
              />
            )}
            <input type="text" name="name" placeholder="Name" />
            <input type="text" name="alias" placeholder="Alias" />
            {/* TODO list of emails */}
            {/* TODO private property to emails */}
            {emails !== undefined &&
              emails.map((email, index) => (
                <div key={index}>
                  {email.main && (
                    <>
                      <label htmlFor={`email-${index}`}>Main Email</label>
                      <br />
                    </>
                  )}
                  <input
                    type="text"
                    name={`email-${index}`}
                    placeholder={
                      index === 0 ? "Email Primary" : `Email ${index + 1}`
                    }
                    defaultValue={email.email}
                  />
                  <br />
                  <label htmlFor={`email-${index}`}>Private</label>
                  <input
                    type="checkbox"
                    defaultChecked={email.private}
                    name={`email-private-${index}`}
                  />
                </div>
              ))}
            {org.$description !== undefined && org.$description === "" && (
              <AlertDismissible
                show
                content={"Please add a description"}
                variant="info"
              />
            )}
            <textarea
              type="text"
              name="description"
              placeholder="Description"
            />
            {org.$address !== undefined && org.$address === "" && (
              <AlertDismissible
                show
                content={`Please add an Address,
                  You can use the map icon to select your address`}
                variant="info"
              />
            )}
            <textarea type="text" name="address" placeholder="Address" />
            <input
              type="text"
              name="location"
              placeholder="Location Lat, Long"
            />
            <div className="mapctrl">
              <Tippy
                arrow={false}
                placement={"right"}
                delay={[1000, 200]}
                content={`${!mapVis ? "Show" : "Hide"} Map`}
              >
                {/* https://github.com/atomiks/tippyjs-react/issues/218 */}
                <i>
                  <SVG
                    className={`mapicon ${!mapVis ? "open" : ""}`}
                    onClick={() => setmapVis(!mapVis)}
                    title={`${!mapVis ? "Show" : "Hide"} Map`}
                    src={mapVis ? svgclose : svgopen}
                  ></SVG>
                </i>
              </Tippy>
            </div>
            <input type="checkbox" name="privateLoc" />
            <input type="checkbox" name="private" />
            <label htmlFor="oldPassword">Change Password</label>
            <input
              type="password"
              name="oldPassword"
              placeholder="Old Password"
            />
            <input type="password" name="newPassword" placeholder="Password" />
            <input
              type="password"
              name="confirm"
              placeholder="Confirm password"
            />
            <label htmlFor="confirm">
              {!passValid &&
                pass &&
                conf &&
                `Passwords don't match ${pass} , ${conf}`}
            </label>
            {clientValidError !== undefined && (
              <AlertDismissible
                show
                content={clientValidError}
                variant="warning"
              />
            )}
            <Button onClick={handleSubmit}>Update</Button>
            {done !== undefined && done && (
              <AlertDismissible
                show
                content={"Updated successfully!"}
                variant="success"
              />
            )}
            {sending !== undefined && sending && (
              <section {...containerProps}>{indicatorEl}</section>
            )}
          </form>
          {mapVis && <Map copyCallback={useCallback} />}
        </div>
      )}
      {serverValidError !== undefined && (
        <AlertDismissible show content={serverValidError} variant="error" />
      )}
    </div>
  );
}

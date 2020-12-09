import { Puff, useLoading } from "@agney/react-loading";
import debounce from "debounce";
import { useEffect, useState } from "react";
import SVG from "react-inlinesvg";
import { Link } from "react-router-dom";
import Org from "../../models/org";
import "./index.css";
import Logout from "./logout";
import svgnot from "./no.svg";
import svgok from "./ok.svg";

export default function Register() {
  const [org] = useState(new Org());
  const [passValid, setPassValid] = useState(false);
  const [pass, setPass] = useState();
  const [conf, setConf] = useState();

  const [loggedin, setLoggedin] = useState();
  const [aliasAvailable, setAliasAvailable] = useState();
  const [aliasAvailableError, setAliasAvailableError] = useState();

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

  const updateOrg = (e) => {
    if (e.target.name.startsWith("email")) {
      // form of email-0, email-1 etc.
      // 0 being primary
      let idx = e.target.name.split("-")[1];
      org["email"](e.target.value, idx);
    } else if (e.target.name === "password" || e.target.name === "confirm") {
      // set name = value
      org[e.target.name](e.target.value);
      validatePass();
    } else {
      // set name = value
      org[e.target.name](e.target.value);
    }
  };

  // https://stackoverflow.com/questions/4220126/run-javascript-function-when-user-finishes-typing-instead-of-on-key-up#comment85608718_16324620
  const checkAliasAvaliable = debounce(
    () => {
      if (org.$alias.length > 3) {
        org.aliasAvailable().then((res) => {
          // possible status codes 500, 200, 403
          switch (res.status) {
            case 200:
              setAliasAvailable(true);
              break;
            case 403:
              setAliasAvailable(false);
              break;
            case 500:
              res.json().then((x) => {
                setAliasAvailableError(x.message);
              });
              break;

            default:
              break;
          }
        });
      } else {
        // TODO show alias > 3 digits message
      }
    },
    400,
    // not immediately
    false
  );

  useEffect(() => {
    org
      .loggedin()
      .then((x) => {
        // https://stackoverflow.com/a/54118576/8608146
        if (x.status !== 200) {
          throw new Error("Not logged in");
        }
        return x.json();
      })
      .then((x) => {
        setLoggedin(true);
        console.log(x);
      })
      .catch(() => {
        setLoggedin(false);
      });
  }, [org]);

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
    if (
      org.$alias.length > 3 &&
      aliasAvailable
    ) {
    } else {
      console.error("Bad alias");
      return;
    }
    org
      .create()
      .then((res) => res.json())
      .then((data) => console.log(data))
      .catch((err) => console.error(err));
  };
  return (
    <div>
      <h2>Register</h2>
      <Link to="/login">Login</Link>
      {loggedin !== undefined ? (
        !loggedin ? (
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
            <input
              type="text"
              name="alias"
              onKeyUp={checkAliasAvaliable}
              placeholder="Alias"
            />
            <label htmlFor="alias">
              {aliasAvailable !== undefined &&
                (aliasAvailable ? (
                  <div>
                    <SVG className="svgicon" title={"Available"} src={svgok}>
                      <div>✅</div>
                    </SVG>
                  </div>
                ) : (
                  <div>
                    <SVG
                      className="svgicon"
                      title={`${org.$alias} Not available`}
                      src={svgnot}
                    >
                      <div>❌</div>
                    </SVG>
                  </div>
                ))}
              {aliasAvailableError !== undefined && aliasAvailableError && (
                <div>{aliasAvailableError}</div>
              )}
            </label>
            {/* TODO list of emails */}
            {/* TODO private property to emails */}
            <input type="text" name="email-0" placeholder="Email" />
            <input type="password" name="password" placeholder="Password" />
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
            <button type="submit">Register</button>
          </form>
        ) : (
          <div>
            You're already loggedin
            <Logout org={org} redirect="/" timeoutDur={5} />
          </div>
        )
      ) : (
        <section {...containerProps}>{indicatorEl}</section>
      )}
    </div>
  );
}

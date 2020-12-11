import { Puff, useLoading } from "@agney/react-loading";
import { useCallback, useEffect, useState } from "react";
import { Button } from "antd";
import { Link, Redirect } from "react-router-dom";
import Org from "../../models/org";
import "./index.css";
import Logout from "./logout";

export default function Login() {
  // Org is the API methods provider as well as store
  const [org] = useState(new Org());

  // Used to track if logged in or not
  const [loggedin, setLoggedin] = useState();

  // if we're done logging in so we can redirect to /dashboard
  const [done, setDone] = useState();

  // tracks if we're sending a request to server to show a loading spinner
  const [sending, setSending] = useState();

  // any error received from the server which will be shown to the user
  const [serverValidityErrors, setServerValidityErrors] = useState();

  // a loading spinner thing
  const { containerProps, indicatorEl: loaderSpinner } = useLoading({
    loading: true,
    indicator: <Puff width="50" />,
  });

  const updateOrg = (e) => {
    if (e.target.name === "email-alias") {
      org["emailOrAlias"](e.target.value);
    } else {
      // set name = value
      org[e.target.name](e.target.value);
    }
  };

  // https://stackoverflow.com/a/53215514/8608146
  // we use it after logout
  const [reload, updateState] = useState();
  const reloadPage = useCallback(() => {
    // we're not logged in anymore
    setLoggedin(undefined);
    updateState({});
  }, []);

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
      .then(() => {
        setLoggedin(true);
      })
      .catch(() => {
        setLoggedin(false);
      });
  }, [org, reload]);

  const handleSubmit = (e) => {
    e.preventDefault();
    setSending(true);
    org
      .login()
      .then(async (res) => {
        setSending(false);
        const jsonD = await res.json();
        switch (res.status) {
          case 422:
            console.error(jsonD["error"]);
            setServerValidityErrors(jsonD["messages"]);
            break;
          case 200:
            // successfully loggedin redirect to dashboard
            setDone(true);
            break;
          case 500:
            setServerValidityErrors([jsonD["error"]]);
            break;
          // login failed StatusForbidden
          case 403:
            setServerValidityErrors([jsonD["messages"]]);
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
      <h2>Login</h2>
      <Link to="/register">Register</Link>
      {loggedin !== undefined ? (
        !loggedin ? (
          <form
            onChange={updateOrg}
            onSubmit={handleSubmit}
            id={"formx"}
            style={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <input
              type="text"
              name="email-alias"
              placeholder="Email or Alias"
            />
            <input type="password" name="password" placeholder="Password" />
            <Button onClick={handleSubmit}>Login</Button>
            {serverValidityErrors !== undefined && (
              <div>
                {serverValidityErrors.map((x, i) => (
                  <p key={i}>{x}</p>
                ))}
              </div>
            )}
            {sending !== undefined && sending && (
              <section {...containerProps}>{loaderSpinner}</section>
            )}
            {done !== undefined && done && (
              <div>
                <Redirect to={"/dashboard"} />
              </div>
            )}
          </form>
        ) : (
          <div>
            You're already loggedin
            <Logout
              org={org}
              redirect="/login"
              timeoutDur={0}
              callback={reloadPage}
            />
          </div>
        )
      ) : (
        <section {...containerProps}>{loaderSpinner}</section>
      )}
    </div>
  );
}

import { Puff, useLoading } from "@agney/react-loading";
import { useCallback, useEffect, useState } from "react";
import { Link, Redirect } from "react-router-dom";
import Org from "../../models/org";
import "./index.css";
import Logout from "./logout";

export default function Login() {
  const [org] = useState(new Org());

  const [loggedin, setLoggedin] = useState();
  const [done, setDone] = useState();
  const [sending, setSending] = useState();
  const [serverValidError, setServerValidError] = useState();

  const { containerProps, indicatorEl } = useLoading({
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
      .then((x) => {
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
            setServerValidError(jsonD["messages"].join("\n"));
            break;
          case 200:
            // successfully loggedin redirect to dashboard
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
      <h2>Login</h2>
      <Link to="/register">Register</Link>
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
              name="email-alias"
              placeholder="Email or Alias"
            />
            <input type="password" name="password" placeholder="Password" />
            <button type="submit">Login</button>
            {serverValidError !== undefined && <div>{serverValidError}</div>}
            {sending !== undefined && sending && (
              <section {...containerProps}>{indicatorEl}</section>
            )}
            {done !== undefined && done && (
              <div>
                <Redirect to={"/dashborad"} />
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
        <section {...containerProps}>{indicatorEl}</section>
      )}
    </div>
  );
}

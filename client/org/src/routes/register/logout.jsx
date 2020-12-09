import Org from "../../models/org";
import { useEffect, useState } from "react";
import { Redirect } from "react-router-dom";

export default function Logout(props) {
  let org;
  let redir;

  // https://stackoverflow.com/a/59283373/8608146
  let stillMounted = { value: false };
  useEffect(() => {
    stillMounted.value = true;
    return () => (stillMounted.value = false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (props) {
    org = props.org;
    if (!org) {
      org = new Org();
    }
    if (props.redirect) {
      redir = true;
    }
  }

  const timeoutDur = props.timeoutDur || 5;
  const [loggedout, setLoggedout] = useState();
  const [redirect] = useState(redir);
  const [waitDone, setWaitDone] = useState(false);
  const [waitT, setWaitT] = useState(timeoutDur);

  const logOut = () => {
    org
      .logout()
      .then((x) => {
        if (x.status !== 202) {
          throw new Error("Failed to logout");
        }
        return x.json();
      })
      .then((x) => {
        setLoggedout(true);
        setTimeout(() => {
          setWaitDone(true);
        }, 1000 * timeoutDur);
        let sec = timeoutDur;
        let handler = setInterval(() => {
          sec -= 1;
          if (!stillMounted.value) {
            clearInterval(handler);
          } else {
            setWaitT(sec);
          }
        }, 1000);
        console.log(x);
      })
      .catch((err) => {
        setLoggedout(false);
        console.error(err);
      });
  };
  return (
    <div>
      {loggedout === undefined ? (
        <button onClick={logOut}>Logout</button>
      ) : !loggedout ? (
        <div>Logout Failed</div>
      ) : (
        <div>
          Logged out successfully
          {redirect && waitDone ? (
            <Redirect to={props.redirect} />
          ) : (
            <p>
              Redirecting to {props.redirect} in {waitT} seconds
            </p>
          )}
        </div>
      )}
      {props.children}
    </div>
  );
}

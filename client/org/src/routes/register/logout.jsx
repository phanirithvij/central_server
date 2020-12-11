import { useEffect, useState } from "react";
import {Button} from "antd";
import { Redirect } from "react-router-dom";
import Org from "../../models/org";

// Must set timeoutDur={0} to redirect immediately
// Default 3 secs count down
export default function Logout(props) {
  let org;
  let redir;
  let callback = props.callback;

  // https://stackoverflow.com/a/59283373/8608146
  let stillMounted = { value: false };
  useEffect(() => {
    stillMounted.value = true;
    return () => (stillMounted.value = false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  org = props.org;
  if (!org) {
    org = new Org();
  }
  if (props.redirect) {
    redir = true;
  }

  const timeoutDur = props.timeoutDur !== undefined ? props.timeoutDur : 3;
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
        if (timeoutDur > 0) {
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
        }
        console.log(x);
        if (timeoutDur === 0) setWaitDone(true);
      })
      .catch((err) => {
        setLoggedout(false);
        console.error(err);
      });
  };

  useEffect(() => {
    waitDone && callback?.();
  }, [waitDone, callback]);

  return (
    <div>
      {loggedout === undefined ? (
        <Button onClick={logOut}>Logout</Button>
      ) : !loggedout ? (
        <div>Logout Failed</div>
      ) : (
        <div>
          {redirect && waitDone ? (
            <>
              <Redirect to={props.redirect} />
            </>
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

import { lazy } from "react";
import { Link, Route } from "react-router-dom";
import Settings from "./settings";
import Profile from "./settings/profile/profile";
import { Button } from "antd";
const Activity = lazy(() => import("./activity"));

function Dashboard({ match }) {
  let { path } = match;
  return (
    <div>
      <Route path={`${path}/activity`} component={Activity} />
      <Route path={`${path}/settings`} component={Settings} />
      <Route path={`${path}/profile`} component={Profile} />
      {/* TODO show following buttons only when url is exactly dashboard */}
    </div>
  );
}
export default Dashboard;

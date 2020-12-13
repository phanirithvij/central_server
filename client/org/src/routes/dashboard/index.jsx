import { lazy } from "react";
import { Route } from "react-router-dom";
import Settings from "./settings";
import Profile from "./settings/profile/profile";

const Activity = lazy(() => import("./activity"));

function Dashboard({ match }) {
  let { path } = match;
  return (
    <div>
      <Route path={`${path}/activity`} component={Activity} />
      <Route path={`${path}/settings`} component={Settings} />
      <Route path={`${path}/profile`} component={Profile} />
    </div>
  );
}
export default Dashboard;

import { lazy } from "react";
import { Link, Route } from "react-router-dom";
import Settings from "./settings";

const Activity = lazy(() => import("./activity"));

function Dashboard({ match }) {
  let { url, path } = match;
  return (
    <div>
      <Link to={`${url}/activity`}>activity</Link>
      <Link to={`${url}/settings`}>settings</Link>
      <h2>Dashboard</h2>
      <Route path={`${path}/activity`} component={Activity} />
      <Route path={`${path}/settings`} component={Settings} />
    </div>
  );
}
export default Dashboard;

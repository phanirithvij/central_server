import { lazy } from "react";
import { Route } from "react-router-dom";
import Settings from "./settings";

const Activity = lazy(() => import("./activity"));

function Dashboard({ match }) {
  let { path } = match;
  return (
    <div>
      <h2>Dashboard</h2>
      <Route path={`${path}/activity`} component={Activity} />
      <Route path={`${path}/settings`} component={Settings} />
    </div>
  );
}
export default Dashboard;

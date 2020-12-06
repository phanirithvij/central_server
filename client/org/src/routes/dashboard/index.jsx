import { lazy } from "react";
import { Link, Route } from "react-router-dom";

const Activity = lazy(() => import("./activity"));

function Dashboard({ match }) {
  let { url, path } = match;
  return (
    <div>
      <Link to={`${url}/activity`}>activity</Link>
      <h2>Dashboard</h2>
      <Route path={`${path}/activity`} component={Activity} />
    </div>
  );
}
export default Dashboard;

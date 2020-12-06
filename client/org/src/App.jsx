import React, { lazy, Suspense } from "react";
import { BrowserRouter as Router, Link, Route, Switch } from "react-router-dom";
import "./App.css";
import ServerBaseURL from "./utils/server";

const Home = lazy(() => import("./routes/home"));
const Dashboard = lazy(() => import("./routes/dashboard"));
const Register = lazy(() => import("./routes/register"));

console.log(ServerBaseURL);
function App() {
  return (
    <div className="App">
      <Router basename={process.env.REACT_APP_BASE_URL}>
        <div>
          <ul>
            <li>
              <Link to="/">Home</Link>
            </li>
            <li>
              <Link to="/register">Register</Link>
            </li>
            <li>
              <Link to="/dashboard">Dashboard</Link>
            </li>
          </ul>

          <hr />
        </div>
        <Suspense fallback={<div>Loading...</div>}>
          <Switch>
            <Route exact path="/">
              <Home />
            </Route>
            <Route path="/register">
              <Register />
            </Route>
            <Route path="/dashboard" component={Dashboard} />
          </Switch>
        </Suspense>
      </Router>
    </div>
  );
}

export default App;

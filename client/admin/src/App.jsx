import React from "react";
import { BrowserRouter as Router, Link, Route, Switch } from "react-router-dom";
import "./App.css";
import { ReactComponent as Logo } from "./logo.svg";
import ServerBaseURL from "./utils/server";

console.log("ServerBaseURL", ServerBaseURL);

function Home() {
  return (
    <div>
      <h2>Home</h2>
    </div>
  );
}

function Register() {
  return (
    <div>
      <h2>Register</h2>
    </div>
  );
}

function Dashboard() {
  return (
    <div>
      <h2>Dashboard</h2>
    </div>
  );
}

function App() {
  return (
    <div className="App">
      <header className="App-header">
        {/* https://create-react-app.dev/docs/adding-custom-environment-variables/ */}
        <Router basename={process.env.REACT_APP_BASE_URL}>
          <div>
            <ul>
              <li>
                <Link to="/">Admin Home</Link>
              </li>
              <li>
                <Link to="/dashboard">Dashboard</Link>
              </li>
              <li>
                <Link to="/register">Register</Link>
              </li>
            </ul>

            <hr />
            <Switch>
              <Route exact path="/">
                <Home />
              </Route>
              <Route path="/register">
                <Register />
              </Route>
              <Route path="/dashboard">
                <Dashboard />
              </Route>
            </Switch>
          </div>
        </Router>
        <Logo className="App-logo" alt="logo" />
      </header>
    </div>
  );
}

export default App;

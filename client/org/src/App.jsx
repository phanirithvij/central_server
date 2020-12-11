import React, { Suspense } from "react";
import { BrowserRouter as Router, Link, Route, Switch } from "react-router-dom";
import "./App.css";
import Dashboard from "./routes/dashboard";
import Home from "./routes/home";
import NavBar from "./components/Nav";
import Register from "./routes/register";
import Login from "./routes/register/login";
import Logout from "./routes/register/logout";
import ServerBaseURL from "./utils/server";
import "antd/dist/antd.css"; // or 'antd/dist/antd.less'

// const Home = lazy(() => import("./routes/home"));
// const Dashboard = lazy(() => import("./routes/dashboard"));
// const Register = lazy(() => import("./routes/register"));
// https://dev.to/iamandrewluca/react-lazy-without-default-export-4b65
// const Login = lazy(() =>
//   import("./routes/register/login").then((module) => ({
//     default: module.Login,
//   }))
// );

console.log("ServerBaseURL", ServerBaseURL);

function App() {
  return (
    <>
      <div className="App">
        <Router basename={process.env.REACT_APP_BASE_URL}>
          <div>
            <div className="nav">
              <NavBar />
            </div>

            {/* Development warning */}
            {process.env.NODE_ENV !== "production" &&
              (window.location.port === "9090" ? (
                <div>Development: Server rendered assets</div>
              ) : (
                <div>Development: React client</div>
              ))}

            <hr />
          </div>
          <Suspense fallback={<div>Loading...</div>}>
            <Switch>
              <Route exact path="/" component={Home} />
              <Route path="/register" component={Register} />
              <Route path="/login" component={Login} />
              <Route
                path="/logout"
                render={() => <Logout redirect={"/login"} />}
              />
              <Route path="/dashboard" component={Dashboard} />
            </Switch>
          </Suspense>
        </Router>
      </div>
    </>
  );
}

export default App;

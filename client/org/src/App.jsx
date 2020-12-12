// const Home = lazy(() => import("./routes/home"));
// const Dashboard = lazy(() => import("./routes/dashboard"));
// const Register = lazy(() => import("./routes/register"));
// https://dev.to/iamandrewluca/react-lazy-without-default-export-4b65
// const Login = lazy(() =>
//   import("./routes/register/login").then((module) => ({
//     default: module.Login,
//   }))
// );
import { Breadcrumb, Grid, Layout, Tag } from "antd";
import "antd/dist/antd.css"; // or 'antd/dist/antd.less'
import React, { Suspense, useEffect, useState } from "react";
import {
  BrowserRouter as Router,
  Link,
  Redirect,
  Route,
  Switch,
  useLocation,
} from "react-router-dom";
import "./App.css";
import AlertDismissible from "./components/Alert";
import SideNavBar, { NavBarComponet } from "./components/Nav";
import Dashboard from "./routes/dashboard";
import Home from "./routes/home";
import Register from "./routes/register";
import Account from "./routes/register/account";
import Login from "./routes/register/login";
import Logout from "./routes/register/logout";
import ServerBaseURL from "./utils/server";
const { useBreakpoint } = Grid;

const { Header, Content } = Layout;

console.log("ServerBaseURL", ServerBaseURL);

function UseBreakpointDemo() {
  const screens = useBreakpoint();
  return (
    <>
      Current break point:{" "}
      {Object.entries(screens)
        .filter((screen) => !!screen[1])
        .map((screen) => (
          <Tag color="blue" key={screen[0]}>
            {screen[0]}
          </Tag>
        ))}
    </>
  );
}

export default function App() {
  const screens = useBreakpoint();
  // TODO sm -> collapse, xs -> hide, md> => show
  const [collapsed, setCollapsed] = useState(screens.sm && !screens.md);
  console.log(collapsed);
  useEffect(() => {
    console.log(screens);
    setCollapsed(screens.sm && !screens.md);
  }, [screens]);
  return (
    <>
      <Router basename={process.env.REACT_APP_BASE_URL}>
        <Layout>
          <Header
            style={{
              position: "fixed",
              width: "100%",
              paddingLeft: !screens.xs ? (collapsed ? 100 : 220) : 0,
              paddingRight: "3vw",
              zIndex: 2,
            }}
          >
            {console.log(screens.xs && !screens.md, collapsed)}
            {!screens.xs && (
              <SideNavBar collapsed={collapsed} setCollapsed={setCollapsed} />
            )}
            <NavBarComponet mode={"horizontal"} />
          </Header>
          <Layout>
            <Layout
              style={{
                padding: "0 24px 24px",
                marginLeft: !screens.xs ? (collapsed ? 80 : 200) : 0,
                minHeight: "100vh",
                // top navbar height
                paddingTop: "66px",
              }}
            >
              <BreadcrumbBar />
              <Content
                className="site-layout-background"
                style={{
                  padding: 24,
                  margin: 0,
                  minHeight: 280,
                }}
              >
                <div className="App">
                  <DevBar />
                  <Suspense
                    fallback={
                      <div style={{ minHeight: "100vh" }}>Loading...</div>
                    }
                  >
                    <Switch>
                      <Route exact path="/" component={Home} />
                      <Route exact path="/account" component={Account} />
                      <Route
                        path="/login"
                        render={() => <Redirect to={"/account/login"} />}
                      />
                      <Route
                        path="/register"
                        render={() => <Redirect to={"/account/register"} />}
                      />
                      <Route path="/account/register" component={Register} />
                      <Route path="/account/login" component={Login} />
                      <Route
                        path="/logout"
                        render={() => <Logout redirect={"/account/login"} />}
                      />
                      <Route path="/dashboard" component={Dashboard} />
                    </Switch>
                  </Suspense>
                </div>
              </Content>
            </Layout>
          </Layout>
        </Layout>
      </Router>
    </>
  );
}

function DevBar() {
  return (
    <>
      {process.env.NODE_ENV !== "production" && (
        <AlertDismissible
          content={
            <>
              <UseBreakpointDemo />
              <p>
                Development
                {window.location.port === "9000"
                  ? "Server Served files"
                  : "Client React"}
              </p>
            </>
          }
          variant="info"
        />
      )}
    </>
  );
}

function BreadcrumbBar() {
  const location = useLocation();
  const parts = location.pathname.split("/");
  return (
    <Breadcrumb style={{ margin: "16px 0" }}>
      {/*
        first will be empty as eg:
        /dashboard/activity => ['','dashboard','activity']
      */}
      {parts.slice(1).map((x, i) => (
        <Breadcrumb.Item key={i}>
          {/* parts still has "" at index 0 so i+2 brecause slice needs i+1 */}
          <Link to={parts.slice(0, i + 2).join("/")}>
            {capitalize(x === "" ? "home" : x)}
          </Link>
        </Breadcrumb.Item>
      ))}
    </Breadcrumb>
  );
}

// https://stackoverflow.com/a/11524251/8608146
// https://stackoverflow.com/a/7224605/8608146
function capitalize(s) {
  return s && s[0].toUpperCase() + s.slice(1);
}

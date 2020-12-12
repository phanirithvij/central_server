// const Home = lazy(() => import("./routes/home"));
// const Dashboard = lazy(() => import("./routes/dashboard"));
// const Register = lazy(() => import("./routes/register"));
// https://dev.to/iamandrewluca/react-lazy-without-default-export-4b65
// const Login = lazy(() =>
//   import("./routes/register/login").then((module) => ({
//     default: module.Login,
//   }))
// );
import { Breadcrumb, Layout } from "antd";
import "antd/dist/antd.css"; // or 'antd/dist/antd.less'
import React, { Suspense, useState } from "react";
import {
  BrowserRouter as Router,
  Link,
  Redirect,
  Route,
  Switch,
  useLocation,
} from "react-router-dom";
import "./App.css";
import NavBar from "./components/Nav";
import Dashboard from "./routes/dashboard";
import Home from "./routes/home";
import Register from "./routes/register";
import Account from "./routes/register/account";
import Login from "./routes/register/login";
import Logout from "./routes/register/logout";
import ServerBaseURL from "./utils/server";

const { Header, Content, Sider } = Layout;

console.log("ServerBaseURL", ServerBaseURL);

export default function App() {
  const [collapsed, setCollapsed] = useState(false);
  const toggleCollapse = () => setCollapsed(!collapsed);

  return (
    <>
      <Router basename={process.env.REACT_APP_BASE_URL}>
        <Layout>
          <Header
            style={{
              position: "fixed",
              width: "100%",
              paddingLeft: !collapsed ? "220px": "100px",
              paddingRight: "3vw",
              zIndex: 2,
            }}
          >
            <NavBar mode={"horizontal"} />
          </Header>
          <Layout>
            <Sider
              // breakpoint="lg"
              // collapsedWidth="0"
              collapsible
              onCollapse={toggleCollapse}
              collapsed={collapsed}
              style={{
                overflow: "auto",
                height: "100vh",
                position: "fixed",
                left: 0,
                width: 200,
                zIndex: 3,
              }}
              className="site-layout-background"
            >
              <NavBar
                mode="inline"
                style={{ height: "100%", borderRight: 0 }}
              />
            </Sider>
            <Layout
              style={{
                padding: "0 24px 24px",
                marginLeft: !collapsed ? 200 : 80,
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

import {
  LaptopOutlined,
  NotificationOutlined,
  UserOutlined,
} from "@ant-design/icons";
// const Home = lazy(() => import("./routes/home"));
// const Dashboard = lazy(() => import("./routes/dashboard"));
// const Register = lazy(() => import("./routes/register"));
// https://dev.to/iamandrewluca/react-lazy-without-default-export-4b65
// const Login = lazy(() =>
//   import("./routes/register/login").then((module) => ({
//     default: module.Login,
//   }))
// );
import { Breadcrumb, Layout, Menu } from "antd";
import "antd/dist/antd.css"; // or 'antd/dist/antd.less'
import React, { Suspense, useState } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import "./App.css";
import NavBar from "./components/Nav";
import Dashboard from "./routes/dashboard";
import Home from "./routes/home";
import Register from "./routes/register";
import Login from "./routes/register/login";
import Logout from "./routes/register/logout";
import ServerBaseURL from "./utils/server";

const { SubMenu } = Menu;
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
              paddingLeft: !collapsed && "220px",
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
                marginLeft: !collapsed ? 200 : null,
                minHeight: "100vh",
                // top navbar height
                paddingTop: "66px",
              }}
            >
              <Breadcrumb style={{ margin: "16px 0" }}>
                <Breadcrumb.Item>Home</Breadcrumb.Item>
                <Breadcrumb.Item>List</Breadcrumb.Item>
                <Breadcrumb.Item>App</Breadcrumb.Item>
              </Breadcrumb>
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
                      <Route path="/register" component={Register} />
                      <Route path="/login" component={Login} />
                      <Route
                        path="/logout"
                        render={() => <Logout redirect={"/login"} />}
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

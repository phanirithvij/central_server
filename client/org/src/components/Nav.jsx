import { Puff, useLoading } from "@agney/react-loading";
import {
  DesktopOutlined,
  HomeOutlined,
  LoginOutlined,
  LogoutOutlined,
  PieChartOutlined,
  ProfileOutlined,
  SettingOutlined,
  UserAddOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Layout, Menu } from "antd";
import React, { useEffect, useState } from "react";
import { Link, useLocation } from "react-router-dom";
import Org from "../models/org";

const { SubMenu } = Menu;
const { Sider } = Layout;

export default function SideNavBar(props) {
  // const collapsed = props.collapsed;
  const [collapsed, setCollapsed] = useState(props.collapsed);

  const toggleCollapse = () => {
    props.setCollapsed?.(!collapsed);
  };

  useEffect(() => {
    setCollapsed(props.collapsed);
  }, [props.collapsed]);

  return (
    <Sider
      // breakpoint="sm"
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
      <NavBarComponet
        style={{ height: "100%", borderRight: 0 }}
        mode={"inline"}
      />
    </Sider>
  );
}

export function NavBarComponet(props) {
  // https://stackoverflow.com/a/60736742/8608146
  const location = useLocation();
  const { containerProps, indicatorEl } = useLoading({
    loading: true,
    indicator: <Puff width="20" />,
  });
  const [loggedin, setLoggedin] = useState();
  const [reload, setReload] = useState(false);

  const [org] = useState(new Org());

  useEffect(() => {
    org
      .loggedin()
      .then((x) => {
        // https://stackoverflow.com/a/54118576/8608146
        if (x.status !== 200) {
          throw new Error("Not logged in");
        }
        return x.json();
      })
      .then((x) => {
        setLoggedin(true);
      })
      .catch(() => {
        setLoggedin(false);
      });
  }, [org, reload]);

  useEffect(() => {
    setReload(reload=>!reload);
  }, [location.pathname]);

  return (
    <Menu
      selectedKeys={[location.pathname]}
      mode={props.mode}
      style={{ height: "100%" }}
      theme="dark"
    >
      {props.mode === "horizontal" && (
        <Menu.Item key="/" icon={<HomeOutlined />}>
          <Link to={"/"}>Home</Link>
        </Menu.Item>
      )}
      {loggedin !== undefined ? (
        loggedin && (
          <SubMenu
            key="/dashboard"
            icon={<DesktopOutlined />}
            title="Dashboard"
          >
            <Menu.Item key="/dashboard" icon={<DesktopOutlined />}>
              <Link to={"/dashboard"}>Home</Link>
            </Menu.Item>
            <Menu.Item key="/dashboard/activity" icon={<PieChartOutlined />}>
              <Link to={"/dashboard/activity"}>Activity</Link>
            </Menu.Item>
            <Menu.Item key="/dashboard/settings" icon={<SettingOutlined />}>
              <Link to={"/dashboard/settings"}>Settings</Link>
            </Menu.Item>
            <Menu.Item key="/dashboard/profile" icon={<ProfileOutlined />}>
              <Link to={"/dashboard/profile"}>Profile</Link>
            </Menu.Item>
          </SubMenu>
        )
      ) : (
        <section {...containerProps}>{indicatorEl}</section>
      )}

      {props.mode === "horizontal" && (
        <SubMenu
          key="/account"
          style={{ float: "right" }}
          icon={<UserOutlined />}
          title="Account"
        >
          {loggedin !== undefined ? (
            !loggedin ? (
              <>
                <Menu.Item key="/account/register" icon={<UserAddOutlined />}>
                  <Link to={"/account/register"}>Register</Link>
                </Menu.Item>
                <Menu.Item key="/account/login" icon={<LoginOutlined />}>
                  <Link to={"/account/login"}>Login</Link>
                </Menu.Item>
              </>
            ) : (
              <Menu.Item key="/logout" icon={<LogoutOutlined />}>
                <Link to={"/logout"}>Logout</Link>
              </Menu.Item>
            )
          ) : (
            <section {...containerProps}>{indicatorEl}</section>
          )}
        </SubMenu>
      )}
    </Menu>
  );
}

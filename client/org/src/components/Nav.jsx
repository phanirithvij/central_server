import {
  DesktopOutlined,
  HomeOutlined,
  LoginOutlined,
  PieChartOutlined,
  ProfileOutlined,
  SettingOutlined,
  UserAddOutlined,
  UserOutlined,
  LogoutOutlined,
} from "@ant-design/icons";
import { Menu } from "antd";
import React from "react";
import { Link, useLocation } from "react-router-dom";

const { SubMenu } = Menu;

export default function NavBar(props) {
  // https://stackoverflow.com/a/60736742/8608146
  const location = useLocation();
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
      <SubMenu key="/dashboard" icon={<DesktopOutlined />} title="Dashboard">
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
      {props.mode === "horizontal" && (
        <SubMenu
          key="/account"
          style={{ float: "right" }}
          icon={<UserOutlined />}
          title="Account"
        >
          <Menu.Item key="/account/register" icon={<UserAddOutlined />}>
            <Link to={"/account/register"}>Register</Link>
          </Menu.Item>
          <Menu.Item key="/account/login" icon={<LoginOutlined />}>
            <Link to={"/account/login"}>Login</Link>
          </Menu.Item>
          <Menu.Item key="/logout" icon={<LogoutOutlined />}>
            <Link to={"/logout"}>Logout</Link>
          </Menu.Item>
        </SubMenu>
      )}
      {/* {props.mode === "horizontal" && process.env.NODE_ENV !== "production" && (
        // https://stackoverflow.com/a/50883195/8608146
        <Menu.Item key="" style={{ float: "right" }}>
          <div>
            {window.location.port === "9090"
              ? "Development: Server rendered assets"
              : "Development: React client"}
          </div>
        </Menu.Item>
      )} */}
    </Menu>
  );
}

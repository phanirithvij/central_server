import {
  DesktopOutlined,
  HomeOutlined,
  LoginOutlined,
  PieChartOutlined,
  SettingOutlined,
  UserAddOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { Menu } from "antd";
import React from "react";
import { Link, useLocation } from "react-router-dom";

const { SubMenu } = Menu;

export default function NavBar(props) {
  const location = useLocation();
  return (
    <div>
      <Menu
        defaultSelectedKeys={[location.pathname]}
        mode={props.mode}
        theme="dark"
      >
        <Menu.Item key="/" icon={<HomeOutlined />}>
          <Link to={"/"}>Home</Link>
        </Menu.Item>
        <SubMenu key="/dashboard" icon={<DesktopOutlined />} title="Dashboard">
          <Menu.Item key="/dashboard/activity" icon={<PieChartOutlined />}>
            <Link to={"/dashboard/activity"}>Activity</Link>
          </Menu.Item>
          <Menu.Item key="/dashboard/settings" icon={<SettingOutlined />}>
            <Link to={"/dashboard/settings"}>Settings</Link>
          </Menu.Item>
          <Menu.Item key="/dashboard/profile" icon={<SettingOutlined />}>
            <Link to={"/dashboard/profile"}>Profile</Link>
          </Menu.Item>
        </SubMenu>
        <SubMenu key="/account" icon={<UserOutlined />} title="Account">
          <Menu.Item key="/register" icon={<UserAddOutlined />}>
            <Link to={"/register"}>Register</Link>
          </Menu.Item>
          <Menu.Item key="/login" icon={<LoginOutlined />}>
            <Link to={"/login"}>Login</Link>
          </Menu.Item>
        </SubMenu>
        {/* Development warning */}
        {process.env.NODE_ENV !== "production" && (
          // 
          <Menu.Item key="" style={{ float: "right" }}>
            <div>
              {window.location.port === "9090"
                ? "Development: Server rendered assets"
                : "Development: React client"}
            </div>
          </Menu.Item>
        )}
      </Menu>
    </div>
  );
}

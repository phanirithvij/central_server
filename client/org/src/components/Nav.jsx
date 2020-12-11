import {
  DesktopOutlined,
  HomeOutlined,
  LoginOutlined,
  PieChartOutlined,
  SettingOutlined,
} from "@ant-design/icons";
import { Menu } from "antd";
import React from "react";
import { Link, useLocation } from "react-router-dom";

const { SubMenu } = Menu;

export default function NavBar() {
  const location = useLocation();
  return (
    <div>
      <Menu
        defaultSelectedKeys={[location.pathname]}
        mode="horizontal"
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
        </SubMenu>
        <Menu.Item key="/register" icon={<LoginOutlined />}>
          <Link to={"/register"}>Register</Link>
        </Menu.Item>
      </Menu>
    </div>
  );
}

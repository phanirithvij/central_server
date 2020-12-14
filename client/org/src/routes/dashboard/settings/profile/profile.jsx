import {
  DatabaseOutlined,
  EditOutlined,
  EllipsisOutlined,
} from "@ant-design/icons";
import Tippy from "@tippyjs/react";
import { Card, Col, Menu, Row, Switch } from "antd";
import { useState } from "react";
import { Link } from "react-router-dom";
import "tippy.js/dist/tippy.css"; // optional
const { SubMenu } = Menu;

const { Meta } = Card;

function ProfileInt() {
  const [loading, setLoading] = useState(false);
  const onChange = (checked) => {
    setLoading(!checked);
  };

  return (
    <>
      <Switch checked={!loading} onChange={onChange} />

      <Card
        loading={loading}
        style={{ width: 300, marginTop: 16 }}
        actions={[
          <DatabaseOutlined key="datasets" />,
          <Tippy
            content={
              "This will not be visible to public, shown only in this page"
            }
          >
            <Link to={`/dashboard/settings`}>
              <EditOutlined key="edit" />
            </Link>
          </Tippy>,
        ]}
      >
        {/* <Skeleton loading={loading} avatar active> */}
        <Meta title="Card title" description="This is the description" />
        {/* </Skeleton> */}
      </Card>
    </>
  );
}

export default function Profile() {
  return (
    <div>
      <Row>
        <Col xs={2} sm={4} md={6} lg={7} xl={8}></Col>
        <Col xs={20} sm={16} md={12} lg={10} xl={8}>
          <ProfileInt />
        </Col>
        <Col xs={2} sm={4} md={6} lg={7} xl={8}></Col>
      </Row>
    </div>
  );
}

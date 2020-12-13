import { Puff, useLoading } from "@agney/react-loading";
import { DatabaseOutlined } from "@ant-design/icons";
import { Alert, Card, Col, Row } from "antd";
import { useEffect, useState } from "react";
import "tippy.js/dist/tippy.css"; // optional
import { PublicListURL } from "../../utils/server";

const { Meta } = Card;

function OrgCard({ info }) {
  const [loading, setLoading] = useState(!info);
  console.log(info);

  useEffect(() => {
    console.log(info);
    setLoading(!info);
  }, [info]);

  return (
    <>
      <Card
        loading={loading}
        style={{ marginTop: 16 }}
        actions={[
          <>
            <a href={`${info.server}/home`}>
              <DatabaseOutlined key="datasets" />
            </a>
          </>,
        ]}
      >
        <Meta title={info?.name} description={info?.description} />
      </Card>
    </>
  );
}

const ARR = new Array(45);

function HomeInt() {
  const [loading, setLoading] = useState();
  const [error, setError] = useState();
  const [cards, setCards] = useState(ARR);

  // Fetch public info
  useEffect(() => {
    setLoading(true);
    fetch(PublicListURL)
      .then(async (res) => {
        let jsonD = await res.json();
        switch (res.status) {
          case 200:
            console.log(jsonD);
            setCards(jsonD);
            break;

          case 404:
            setError({ ...jsonD, code: 404 });
            break;

          default:
            break;
        }
        setLoading(false);
      })
      .catch((err) => {
        setLoading(false);
        setError({ messages: ['Server could not be reached'], error: "Failed to fetch", code: '' });
      });
  }, []);

  // a loading spinner thing
  const { containerProps, indicatorEl: loaderSpinner } = useLoading({
    loading: true,
    indicator: <Puff width="50" />,
  });

  return (
    <div>
      <Row>
        <Col xs={2} sm={4} md={6} lg={7} xl={8}></Col>
        <Col xs={20} sm={16} md={12} lg={10} xl={8}>
          {error !== undefined && (
            <Alert
              message={`${error.code} ${error.error}`}
              description={error.messages.toString()}
              type={"error"}
            />
          )}
          {cards.map((card, index) => (
            <OrgCard key={index} info={card} />
          ))}
          {loading && <section {...containerProps}>{loaderSpinner}</section>}
        </Col>
        <Col xs={2} sm={4} md={6} lg={7} xl={8}></Col>
      </Row>
    </div>
  );
}

export default function Home() {
  return (
    <div>
      <h2>Home</h2>
      <HomeInt />
    </div>
  );
}

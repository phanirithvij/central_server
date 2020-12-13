import { Puff, useLoading } from "@agney/react-loading";
import { MinusCircleOutlined, PlusOutlined } from "@ant-design/icons";
import Tippy from "@tippyjs/react";
import { Alert, Button, Checkbox, Col, Divider, Form, Input, Row } from "antd";
import { lazy, useEffect, useState } from "react";
import SVG from "react-inlinesvg";
import { Link } from "react-router-dom";
import "tippy.js/dist/tippy.css"; // optional
import AlertDismissible from "../../../components/Alert";
import Org from "../../../models/org";
import Logout from "../../register/logout";
import "./index.css";
import svgopen from "./map.svg";
import svgclose from "./mapclose.svg";

const { TextArea } = Input;
// Must be lazy for it is ~ 1MB, gziped 200 KB
const Map = lazy(() => import("../../../components/Map"));

export default function Settings() {
  const [mapVis, setmapVis] = useState(false);
  const [passValid, setPassValid] = useState(false);
  const [form] = Form.useForm();
  window.form = form;
  const [pass, setPass] = useState();
  const [conf, setConf] = useState();
  const [org] = useState(new Org());
  const [clientValidError, setClientValidError] = useState();
  const [serverValidError, setServerValidError] = useState();
  const [sending, setSending] = useState();
  const [loggedin, setLoggedin] = useState();
  const [loggedinJson, setLoggedinJson] = useState();
  const [done, setDone] = useState();
  const [, setReload] = useState(false);

  const { containerProps, indicatorEl } = useLoading({
    loading: true,
    indicator: <Puff width="50" />,
  });

  const validatePass = () => {
    setPass(org.$password);
    setConf(org._confirm);
    if (!org.$password.startsWith(org._confirm)) {
      // entered wrong thing
      setPassValid(false);
      // console.log(org._confirm, "doesn't match", org.$password);
      return;
    } else {
      if (
        org._confirm.length >= org.$password.length &&
        org._confirm !== org.$password
      ) {
        setPassValid(false);
        // console.log(org._confirm, "doesn't match", org.$password);
        return;
      }
    }
    setPassValid(true);
  };

  useEffect(() => {
    const json = loggedinJson;
    if (json === undefined) return;
    Object.keys(json).forEach((key) => {
      // console.log(key);
      if (key === "id") return;
      if (key === "emails") {
        org["emails"](json[key]);
        // console.log(org._emailList());
        let ems = org._emailList();
        // TODO only one of these
        form.setFieldsValue({ emails: ems });
        form.setFields([{ name: "emails", value: ems }]);
        return;
      }

      // console.log(org[key], key);
      // org[key](json[key]);
      if (typeof json[key] === "boolean") {
        form.setFields([{ name: key, value: json[key] }]);
        // form.setFields([{ name: key, checked: json[key] }]);
        // document.querySelector(`input[name="${key}"]`).checked = json[key];
      } else {
        let input = document.querySelector(`input[name="${key}"]`);
        if (["address", "description"].includes(key)) {
          input = document.querySelector(`textarea[name="${key}"]`);
        }
        form.setFields([{ name: key, value: json[key] }]);
        // input.value = json[key];
        if (input?.disabled) {
          input.placeholder = `Alias:\t${json[key]}`;
        }
        // console.log(key, json[key], input.value);
      }
      // need to set state to refesh any alerts guiding empty fields
      setReload((reload) => !reload);
    });
  }, [loggedinJson, org, form]);

  useEffect(() => {
    org
      .fetch()
      .then(async (x) => {
        switch (x.status) {
          case 401:
            setServerValidError("Not logged in, please login");
            break;

          case 404:
            // Show user that warning and a logout button
            setServerValidError(
              "Organization doesn't exist on our server please mail us if you believe this to be a mistake."
            );
            break;

          case 200:
            // Logged in and got the org details
            // Fill up org
            let json = await x.json();
            setLoggedin(true);
            setLoggedinJson(json);
            break;

          default:
            break;
        }
      })
      .catch((err) => {
        setServerValidError(err.message);
        console.error(err);
      });
  }, [org]);

  // triggered when clicked on the copy icons in the marker popup
  const useCallback = (type, value) => {
    org[type](value);
    form.setFields([{ name: type, value: value }]);
  };

  const handleSubmit = (e) => {
    // clear any previous server errors
    setServerValidError(undefined);
    e.preventDefault();
    const values = form.getFieldsValue();
    if (values.password && values.confirm) {
      // we modfied the password
      if (passValid && values.password === values.confirm) {
      } else {
        setClientValidError("Password is not a valid password");
        return;
      }
    }
    if (values.location) {
    } else {
      setClientValidError("Location is required");
      return;
    }
    // if string convert to array of lat, long
    if (typeof values.location === "string") {
      let loc = values.location.split(",");
      if (loc.length !== 2) {
        setClientValidError("Location is not valid");
        return;
      }
      loc = loc.map((l) => parseFloat(l));
      if (isNaN(loc[0]) || isNaN(loc[1])) {
        setClientValidError("Location is not valid");
        return;
      }
      values.location = loc;
    }
    setSending(true);
    org
      .updateSettings(values)
      .then(async (res) => {
        setSending(false);
        const jsonD = await res.json();
        switch (res.status) {
          case 422:
            console.error(jsonD["error"]);
            setServerValidError(jsonD["messages"].join("\n"));
            break;
          case 200:
            // successfully updted
            // TODO update form
            console.log("Update form???", jsonD);
            setLoggedinJson(jsonD);
            setDone(true);

            // close info in 3 secs
            setTimeout(() => {
              setDone(undefined);
            }, 3000);

            break;
          case 500:
            setServerValidError(jsonD["error"]);
            break;
          default:
            break;
        }
      })
      .catch((err) => {
        setSending(false);
        setServerValidError("Could not reach the server, please try again.");
        console.error(err);
      });
  };
  return (
    <>
      <Row>
        <Col xs={2} sm={4} md={6} lg={7} xl={8}></Col>
        <Col xs={20} sm={16} md={12} lg={10} xl={8}>
          <div>
            {/* Org alias will be undefined only if not loggedin  */}
            {loggedin !== undefined && loggedin ? (
              <>
                <Form form={form} id="formx" style={{ maxWidth: "600px" }}>
                  {form.getFieldValue("name") === "" && (
                    <AlertDismissible
                      show
                      content={"Please add all the details"}
                      variant="info"
                    />
                  )}
                  <Form.Item
                    name="name"
                    rules={[
                      {
                        required: true,
                        message: "Please set a Name",
                      },
                    ]}
                  >
                    <Input name="name" placeholder="Organization Name" />
                  </Form.Item>
                  <Form.Item
                    name="alias"
                    rules={[
                      {
                        required: false,
                        message: "This can't be changed",
                      },
                    ]}
                  >
                    <Input name="alias" disabled placeholder="Alias" />
                  </Form.Item>
                  <Form.Item
                    valuePropName="checked"
                    name="private"
                    label="Organization is private"
                  >
                    <Checkbox name="private" />
                  </Form.Item>
                  {/* TODO list of emails */}
                  {/* TODO private property to emails */}

                  <Divider>Emails</Divider>
                  <Form.List name="emails">
                    {(fields, { add, remove }, { errors }) => (
                      <>
                        {fields.map((field, index) => (
                          <Form.Item key={index}>
                            {form.getFieldValue([
                              "emails",
                              field.name,
                              "main",
                            ]) && <Alert description={"Primary Email"} />}
                            <Row justify="space-around" align="middle">
                              <Col
                                style={{
                                  flex: "0 0 100%",
                                  maxWidth: "85%",
                                }}
                              >
                                <Form.Item
                                  name={[field.name, "email"]}
                                  label="Email"
                                  rules={[{ required: true }]}
                                >
                                  <Input placeholder={`Email ${index}`} />
                                </Form.Item>
                                <Form.Item
                                  label="Private"
                                  name={[field.name, "private"]}
                                  rules={[{ required: true }]}
                                  valuePropName="checked"
                                >
                                  <Checkbox placeholder="Private" />
                                </Form.Item>
                                <Form.Item
                                  name={[field.name, "main"]}
                                  rules={[{ required: true }]}
                                  hidden
                                >
                                  <Input />
                                </Form.Item>
                                <Form.Item
                                  name={[field.name, "id"]}
                                  rules={[{ required: true }]}
                                  hidden
                                >
                                  <Input />
                                </Form.Item>
                              </Col>
                              <Col span={1}>
                                {!form.getFieldValue([
                                  "emails",
                                  field.name,
                                  "main",
                                ]) && (
                                  <MinusCircleOutlined
                                    className="dynamic-delete-button"
                                    onClick={() => remove(field.name)}
                                  />
                                )}
                              </Col>
                            </Row>
                          </Form.Item>
                        ))}
                        <Form.Item>
                          <Button
                            type="dashed"
                            onClick={() => add()}
                            style={{ width: "60%" }}
                            icon={<PlusOutlined />}
                          >
                            Add Email
                          </Button>
                        </Form.Item>
                      </>
                    )}
                  </Form.List>

                  {form.getFieldValue("description") === "" && (
                    <AlertDismissible
                      show
                      content={"Please add a description"}
                      variant="info"
                    />
                  )}
                  <Form.Item
                    name="description"
                    rules={[
                      {
                        required: true,
                        message: "Please add Description",
                      },
                    ]}
                  >
                    <TextArea
                      autoSize
                      showCount
                      maxLength={600}
                      name="description"
                      placeholder="Description"
                    />
                  </Form.Item>

                  <Divider orientation="left">Address and Location</Divider>

                  {form.getFieldValue("address") === "" && (
                    <AlertDismissible
                      show
                      content={`Please add an Address,
                  You can use the map icon to select your address`}
                      variant="info"
                    />
                  )}
                  <Form.Item name="address">
                    <TextArea autoSize name="address" placeholder="Address" />
                  </Form.Item>
                  <Form.Item
                    name="location"
                    rules={[
                      {
                        required: true,
                        message: "Please input Location",
                      },
                    ]}
                  >
                    <Input name="location" placeholder="Location Lat, Long" />
                  </Form.Item>
                  <div className="mapctrl">
                    <Tippy
                      arrow={false}
                      placement={"right"}
                      delay={[1000, 200]}
                      content={`${!mapVis ? "Show" : "Hide"} Map`}
                    >
                      {/* https://github.com/atomiks/tippyjs-react/issues/218 */}
                      <i>
                        <SVG
                          className={`mapicon ${!mapVis ? "open" : ""}`}
                          onClick={() => setmapVis(!mapVis)}
                          title={`${!mapVis ? "Show" : "Hide"} Map`}
                          src={mapVis ? svgclose : svgopen}
                        ></SVG>
                      </i>
                    </Tippy>
                  </div>
                  {/* https://github.com/ant-design/ant-design/issues/20803#issuecomment-601626759 */}
                  <Form.Item
                    valuePropName="checked"
                    name="privateLoc"
                    label="Location is private"
                  >
                    <Checkbox name="privateLoc" />
                  </Form.Item>

                  <Divider orientation="left">Server Configuration</Divider>
                  <Form.Item name="serverAlias">
                    <Input name="serverAlias" placeholder="Server nick name" />
                  </Form.Item>
                  <Form.Item name="server">
                    <Input name="server" placeholder="Server URL" />
                  </Form.Item>
                  <Form.Item
                    name={"serverID"}
                    rules={[{ required: true }]}
                    hidden
                  >
                    <Input />
                  </Form.Item>

                  <Divider orientation="left">Change Password</Divider>
                  <Form.Item name="oldPassword">
                    <Input.Password
                      name="oldPassword"
                      placeholder="Old Password"
                    />
                  </Form.Item>
                  <Form.Item name="newPassword">
                    <Input.Password name="newPassword" placeholder="Password" />
                  </Form.Item>
                  <Form.Item name="confirm">
                    <Input.Password
                      type="password"
                      name="confirm"
                      placeholder="Confirm password"
                    />
                  </Form.Item>
                  <label htmlFor="confirm">
                    {!passValid &&
                      pass &&
                      conf &&
                      `Passwords don't match ${pass} , ${conf}`}
                  </label>
                  {clientValidError !== undefined && (
                    <AlertDismissible
                      show
                      content={clientValidError}
                      variant="warning"
                    />
                  )}
                  <Button onClick={handleSubmit}>Update</Button>
                </Form>
                {done !== undefined && done && (
                  <AlertDismissible
                    show
                    content={"Updated successfully!"}
                    variant="success"
                  />
                )}
                {sending !== undefined && sending && (
                  <section {...containerProps}>{indicatorEl}</section>
                )}
                {mapVis && <Map copyCallback={useCallback} />}
              </>
            ) : (
              <section {...containerProps}>{indicatorEl}</section>
            )}
            {serverValidError !== undefined && (
              <AlertDismissible
                show
                header={serverValidError}
                content={
                  <>
                    {serverValidError.includes("login") && (
                      <Button>
                        <Link to={"/account/login"}>Login</Link>
                      </Button>
                    )}
                    {serverValidError.includes("doesn't exist") && (
                      <Logout
                        org={org}
                        redirect="/account/login"
                        timeoutDur={0}
                      />
                    )}
                  </>
                }
                variant="error"
              />
            )}
          </div>
        </Col>
        <Col xs={2} sm={4} md={6} lg={7} xl={8}></Col>
      </Row>
    </>
  );
}

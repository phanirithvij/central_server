import { Puff, useLoading } from "@agney/react-loading";
import { UserOutlined } from "@ant-design/icons";
import { Button, Col, Form, Input, Row } from "antd";
import debounce from "debounce";
import { useEffect, useState } from "react";
import SVG from "react-inlinesvg";
import { Link, Redirect } from "react-router-dom";
import AlertDismissible from "../../components/Alert";
import Org from "../../models/org";
import "./index.css";
import Logout from "./logout";
import svgnot from "./no.svg";
import svgok from "./ok.svg";

export default function Register() {
  const [org] = useState(new Org());
  const [form] = Form.useForm();
  const [passValid, setPassValid] = useState(false);
  const [pass, setPass] = useState();
  const [conf, setConf] = useState();

  const [loggedin, setLoggedin] = useState();
  const [aliasAvailable, setAliasAvailable] = useState();
  const [aliasAvailableError, setAliasAvailableError] = useState();
  const [clientValidError, setClientValidError] = useState();
  const [serverValidError, setServerValidError] = useState();
  const [done, setDone] = useState();
  const [sending, setSending] = useState();

  const { containerProps, indicatorEl } = useLoading({
    loading: true,
    indicator: <Puff width="50" />,
  });

  const validatePass = () => {
    setPass(org.$password);
    setConf(org._confirm);
    if (org.$password === undefined) {
      setPassValid(false);
      return;
    }
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

  const updateOrg = (e) => {
    if (e.target.name.startsWith("email")) {
      // form of email-0, email-1 etc.
      // 0 being primary
      let idx = e.target.name.split("-")[1];
      org["email"]({ email: e.target.value }, idx);
    } else if (e.target.name === "password" || e.target.name === "confirm") {
      // set name = value
      org[e.target.name](e.target.value);
      validatePass();
    } else {
      // set name = value
      org[e.target.name](e.target.value);
    }
  };

  // https://stackoverflow.com/questions/4220126/run-javascript-function-when-user-finishes-typing-instead-of-on-key-up#comment85608718_16324620
  const checkAliasAvaliable = debounce(
    () => {
      if (org.$alias && org.$alias.length > 3) {
        org.aliasAvailable().then((res) => {
          // possible status codes 500, 200, 403
          switch (res.status) {
            case 200:
              setAliasAvailable(true);
              setAliasAvailableError(undefined);
              break;
            case 403:
              setAliasAvailable(false);
              setAliasAvailableError(undefined);
              break;
            case 500:
              res.json().then((x) => {
                setAliasAvailableError(x.message);
              });
              break;

            default:
              break;
          }
        });
      } else {
        // TODO show alias > 3 digits message
        setAliasAvailable(undefined);
        setAliasAvailableError("Alias must be 4 or more characters long");
      }
    },
    400,
    // not immediately
    false
  );

  let stillMounted = { value: false };
  useEffect(() => {
    stillMounted.value = true;
    return () => (stillMounted.value = false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (stillMounted.value)
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
          if (stillMounted.value) setLoggedin(true);
        })
        .catch(() => {
          if (stillMounted.value) setLoggedin(false);
        });
  }, [org, stillMounted.value]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (
      passValid &&
      org.$password &&
      org._confirm &&
      org.$password === org._confirm
    ) {
    } else {
      setClientValidError("Password is not a valid password");
      return;
    }
    if (!aliasAvailable) {
      setClientValidError("Alias is not available");
      return;
    }
    setSending(true);
    org
      .create()
      .then(async (res) => {
        setSending(false);
        const jsonD = await res.json();
        switch (res.status) {
          case 422:
            console.error(jsonD["error"]);
            setServerValidError(jsonD["messages"].join("\n"));
            break;
          case 201:
            // successfully created org redirect to dashboard
            setDone(true);
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
        console.error(err);
      });
  };
  return (
    <>
      <Row>
        <Col xs={2} sm={4} md={6} lg={7} xl={8}></Col>
        <Col xs={20} sm={16} md={12} lg={10} xl={8}>
          <div>
            <h2>Register</h2>
            <Link to="/account/login">Login</Link>
            {loggedin !== undefined ? (
              !loggedin ? (
                <>
                  <Form
                    form={form}
                    onChange={updateOrg}
                    id="formx"
                    style={{ maxWidth: "600px" }}
                  >
                    <Form.Item>
                      <Input
                        placeholder="Alias"
                        onKeyUp={checkAliasAvaliable}
                        name="alias"
                      />
                      <label htmlFor="alias">
                        {aliasAvailable !== undefined && (
                          <div
                            style={{
                              display: "flex",
                              justifyContent: "center",
                              alignItems: "center",
                            }}
                          >
                            <SVG
                              className="svgicon"
                              title={"Available"}
                              src={aliasAvailable ? svgok : svgnot}
                            >
                              <div>{aliasAvailable ? "✅" : "❌"} </div>
                            </SVG>
                            {org.$alias} is{" "}
                            {aliasAvailable ? "Available" : "Taken"}
                          </div>
                        )}
                        {aliasAvailableError !== undefined &&
                          aliasAvailableError && (
                            <AlertDismissible
                              show
                              content={aliasAvailableError}
                              variant="error"
                            />
                          )}
                      </label>
                    </Form.Item>
                    <Form.Item>
                      <Input
                        name="email-0"
                        prefix={<UserOutlined />}
                        placeholder="Primary Email"
                      />
                    </Form.Item>
                    <Form.Item
                      rules={[
                        {
                          required: true,
                          message: "Please input your password!",
                        },
                      ]}
                    >
                      <Input.Password placeholder="Password" name="password" />
                    </Form.Item>
                    <Form.Item
                      rules={[
                        {
                          required: true,
                          message: "Please input your password!",
                        },
                      ]}
                    >
                      <Input.Password
                        placeholder="Confirm Password"
                        name="confirm"
                      />
                    </Form.Item>
                    <label htmlFor="confirm">
                      {!passValid &&
                        pass &&
                        conf &&
                        `Passwords don't match ${pass} , ${conf}`}
                    </label>
                    <Form.Item>
                      <Button type="primary" onClick={handleSubmit}>
                        Register
                      </Button>
                    </Form.Item>
                  </Form>
                  {clientValidError !== undefined && (
                    <AlertDismissible show content={clientValidError} />
                  )}
                  {serverValidError !== undefined && (
                    <div>{serverValidError}</div>
                  )}
                  {sending !== undefined && sending && (
                    <section {...containerProps}>{indicatorEl}</section>
                  )}
                  {done !== undefined && done && (
                    <div>
                      <Redirect to={"/dashboard/settings"} />
                    </div>
                  )}
                </>
              ) : (
                <AlertDismissible
                  variant="info"
                  content={
                    <>
                      You're already loggedin
                      <Logout
                        org={org}
                        redirect="/account/login"
                        timeoutDur={0}
                      />
                    </>
                  }
                />
              )
            ) : (
              <section {...containerProps}>{indicatorEl}</section>
            )}
          </div>
        </Col>
        <Col xs={2} sm={4} md={6} lg={7} xl={8}></Col>
      </Row>
    </>
  );
}

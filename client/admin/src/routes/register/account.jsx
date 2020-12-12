import { Button } from "antd";
import { Link } from "react-router-dom";

export default function Account() {
  return (
    <>
      <h3>Account</h3>
      <Button type="primary">
        <Link to="/account/register">Register</Link>
      </Button>
      <Button>
        <Link to="/account/login">Login</Link>
      </Button>
    </>
  );
}

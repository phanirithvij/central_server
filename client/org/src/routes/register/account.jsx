import { Link } from "react-router-dom";

export default function Account() {
  return (
    <>
      <h3>Account</h3>
      <Link to="/account/register">Register</Link>
      <Link to="/account/login">Login</Link>
    </>
  );
}

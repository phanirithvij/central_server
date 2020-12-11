import { Alert } from "antd";

export default function AlertDismissible(props) {
  return (
    <Alert
      message="Error Text"
      description="Error Description Error Description Error Description Error Description Error Description Error Description"
      type="error"
      closable
    >
      {props.content && <div>{props.content}</div>}
    </Alert>
  );
}

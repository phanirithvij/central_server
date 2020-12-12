import { Alert } from "antd";

export default function AlertDismissible(props) {
  return (
    <Alert
      message={props.header}
      description={props.content}
      type={props.variant}
      closable
    >
    </Alert>
  );
}

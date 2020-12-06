// in prod server is /, in dev server is :9090
let ServerBaseURL =
  process.env.NODE_ENV === "production" ? "/" : "http://localhost:9090/";

export default ServerBaseURL;

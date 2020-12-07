// in prod server is /, in dev server is :9090
const ServerBaseURL =
  process.env.NODE_ENV === "production" ? "/" : "http://localhost:9090/";
export const RegisterURL = ServerBaseURL + "register/";
export const OrgSettingsURL = ServerBaseURL + "settings/";

export default ServerBaseURL;
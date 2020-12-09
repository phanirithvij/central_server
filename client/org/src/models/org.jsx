import { OrgSettingsURL, RegisterURL, OrgAliasCheckURL, LogoutURL } from "../utils/server";
export default class Org {
  constructor(props) {
    this.props = props;
  }
  /**
   * @param {string} addr
   */
  address(addr) {
    this.$address = addr;
    return this;
  }
  /**
   * @param {[number, number]} loc
   */
  location(loc) {
    this.$location = loc;
    return this;
  }
  /**
   * @param {{email: string, private: boolean}} em
   * @param {number} indx
   */
  email(em, indx) {
    if (!this.$emails) this.$emails = {};
    this.$emails[indx] = em;
    return this;
  }
  /**
   * @param {string} n
   */
  name(n) {
    this.$name = n;
    return this;
  }
  /**
   * @param {string} a
   */
  alias(a) {
    this.$alias = a;
    return this;
  }
  /**
   * @param {string} d
   */
  description(d) {
    this.$description = d;
    return this;
  }
  /**
   * @param {string} p
   */
  password(p) {
    this.$password = p;
    return this;
  }
  /**
   * @param {boolean} b
   */
  private(b) {
    this.$private = b;
    return this;
  }
  /**
   * @param {string} p
   */
  confirm(p) {
    this._confirm = p;
    return this;
  }
  /**
   * @param {boolean} update
   */
  create(update) {
    if (!this.props) this.props = {};
    Object.keys(this).forEach((k) => {
      if (k.startsWith("$")) {
        if (k === "$emails") {
          this.props[k.split("$")[1]] = Object.values(this[k]);
        } else this.props[k.split("$")[1]] = this[k];
      }
    });

    const url = update ? OrgSettingsURL : RegisterURL;
    return fetch(url, {
      method: update ? "PUT" : "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(this.props),
    });
  }
  // Call this from the settings route to update the Org
  update() {
    return this.create(true);
  }

  // fetch the org from the server
  fetch() {
    return fetch(OrgSettingsURL, {
      method: "GET",
      credentials: "include",
    });
  }

  loggedin() {
    return fetch(RegisterURL, {
      method: "GET",
      credentials: "include",
    });
  }
  logout() {
    return fetch(LogoutURL, {
      method: "GET",
      credentials: "include",
    });
  }
  aliasAvailable() {
    // const url = new URL(OrgAliasCheckURL);
    // url.searchParams.set("alias", this.$alias);
    return fetch(`${OrgAliasCheckURL}/${this.$alias}`, {
      method: "GET",
    });
  }
}

window.Org = Org;

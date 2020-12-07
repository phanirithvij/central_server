import { OrgSettingsURL, RegisterURL } from "../utils/server";
export default class Org {
  constructor(props) {
    this.props = props;
  }
  /**
   * @param {string} addr
   */
  address(addr) {
    this._address = addr;
    return this;
  }
  /**
   * @param {[number, number]} loc
   */
  location(loc) {
    this._location = loc;
    return this;
  }
  /**
   * @param {string} em
   * @param {number} indx
   */
  email(em, indx) {
    if (!this._emails) this._emails = [];
    this._emails[indx] = em;
    return this;
  }
  /**
   * @param {string} n
   */
  name(n) {
    this._name = n;
    return this;
  }
  /**
   * @param {string} a
   */
  alias(a) {
    this._alias = a;
    return this;
  }
  /**
   * @param {string} d
   */
  description(d) {
    this._description = d;
    return this;
  }
  /**
   * @param {string} p
   */
  password(p) {
    this._password = p;
    return this;
  }
  confirm(p) {
    this._confirm = p;
    return this;
  }
  /**
   *
   * @param {boolean} update
   */
  create(update) {
    this.props = {
      password: this._password,
      emails: this._emails,
      location: this._location,
      name: this._name,
      description: this._description,
      alias: this._alias,
      address: this._address,
    };
    const url = update ? OrgSettingsURL : RegisterURL;
    return fetch(url, {
      method: update ? "PUT" : "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(this.props),
    });
  }
  update() {
    return this.create(true);
  }
}

window.Org = Org;

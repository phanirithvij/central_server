# ARCHIVE

Could've been made to be more useful if not for time constraints.
Tenants implemented having never heard of the concept, so it is bad.

# Central Server

For technical and implementation details read [./README.tech.md](./README.tech.md)

> This is a Work In Progress (WIP) and there will definitely be breaking, business logic changes

This is the code for the central organization meta managing server which is required or optional based on the decision that's yet to be made at the [initial specification](https://docs.google.com/document/d/1KT_sxSCRaaK5xkVE1t0xDTUzlPYNQJAJXs4B3K5YDkM/view)

## Notes for users (organizations)

If **Optional**

- The central database becomes optional and the corpora server can be used independantly but the client needs to be configured accordingly
- The client needs to be a custom client different from the centralized client which is typically found in the App / Play stores. And the client needs to be distributed to the users on your own.

If **Required**

- The central database details needs to be configured on the server side.

  - Go to the central organization dashboard TODO:URL and register a new organization and configure your server to the provided credentials and instructions by the dashboard.
  - Or after creating your server profile you can register from inside the server config dashboard

- For public organizations
  - Users of the clients will see your organization once you add in data after a grace period of 4-5 days (undecided) after manual verification by the central orgainzation and manual approval.
- For private organizations which intend to only collect data but not reveal themselves to the users (of client apps), their corpora data will be shown when they don't choose to go in any particular organization's page and in the home page

## TODO:

- Applicable to all the organizations which register themselves
  - Privacy policy for the central dashboard
  - Terms of service for the central dashboard
- Decide on the Licenses of the codebases for
  - server
  - client
  - central server
- Decide on the valid opensource corpora usage licenses for public organizations
- Organization corpora requests from one organization to other (clarify?)
- For private users
    - To resolve any violation of private user voice samples we need to maintain the organizations which have access to the private recording so we need to maintain a database for that.
  - Need to mention this in the data collection section for the Privacy Policy

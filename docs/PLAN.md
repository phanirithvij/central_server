Organization

- Private bool
- Name string
- Email[] - emails for hub/user communication
- ID string serverAssigned
- Alias string - Human friendly org slug serverRecommended
- Description string - human readable description
- LocationStr string - Manual location address
- Location
- Server (servers ??)

Email

- Email string
- Private bool - If it can not be shown to public
- Primary bool

Location - Pick on map or something

- Longitude string
- Latitude string
- Private - If it can not to be shown on profile

Server

- URL URL
- Alias string
- Timezone string
- IsOnline? bool
- IsVerified? bool
- BanDetails
  - Banned bool
  - Reason string - Why we banned it?
- DeletionDetails
  - reason string - Why it was deleted?
  - DeletedOn time.Time - When it was deleted
- Suspensions
  - Incidents() - SuspensionIncident[] - Gets suspension incidents from DB
- StatusReason string - Reason for ban/suspension
- Status ServerStatus:enum - [Banned|Suspended|Online|Offline|Deleted]
- Details ??

SuspensionIncident

- Reason string - Why we suspended it?
- DurationReason string - Why we suspend it for this duration?
- SuspendedTill time.Time - Suspended till this datetime
- SuspendedOn time.Time - Suspended on this datetime

ServerStatus

- Banned - We banned it
- Suspended - We suspended it
- Online - Server is online
- Offline - Server is offline
- Deleted - Server was deleted by the organization

Admin

- Username string
- Name string
- PasswordHash string
- Email string - Admin's email address
- Main bool - If admin is the main admin
- AddedBy string - Username of admin who added this admin
- TimeZone string
- Capabilites[]
- Orgs
  - CanAddOrgs()
  - CanEditOrgs()
  - CanRemoveOrgs()
  - CanEditAddedOrgs()
  - CanDeleteAddedOrgs()
- Admins
  - CanAddAdmin()
  - CanEditAdmin()
  - CanEditAddedAdmins()
  - CanRemoveAddedAdmins()

Capability ??

- Type //Activity|Logs|Analytics - Admin can view the page
  - To prevent some admins exporting everything
- Allowed bool
- By string //Username of admin who changed this
- UpdatedOn time.Time //Added/Updated on this datetime

Hub

- Admins[] - list of admins managing the hub
- SuspensionDefaultDur - time.Duration - 3 days

**IF TIME PERMITS**

Status.Hub.org

- Status string - Online|Offline
- Incidents[] - All the incidents in the past

Incident

- ID string
- MainReason string //Reason for the Incident
- Timeline[]
  - Description string //html or markdown string
  - Status //Resolved|Pending
  - Time time.Time //Current time

## DEADLINE PLAN

> TODAY

- [ ] Admin ui
  - [ ] single admin
  - [ ] simple auth
- [ ] Dashboard
  - [ ] Orgs
  - [ ] Delete or Suspend orgs
    - [ ] Just gorm delete
    - [ ] Can undelete if wanted
    - [ ] Delete reason -> [Deletd by Organization, Banned, Suspended]

- [ ] Expose Org Server Homepage so we can link to it
  - [x] At $serverURL/home
  - [ ] (Later) Open Link in app
    - https://pub.dev/packages/uni_links

- [ ] Org readonly activity dashboard
  - [ ] Any update changes in settings
- [x] Org dashboard profile
- [x] List of public orgs on homepage
  - [x] Api enpoint

- [ ] Integrate filebrowser from fate
  - [ ] Integrate filebrowser like fate but don't use fate
- [ ] Integrate fate into org server and central both
  - [ ] Filebrowser branding subprocess call in python for org server

> TOMMOROW

- [ ] Org server url connect to central server
- [ ] Dynamic DB split and serve to clients
- [ ] Client can login via central server
  - [ ] Server can check if client is authenticated in central server
    - [ ] If not show the org server's login page
      - [ ] Send these credentials to central server

- [ ] Client sees list of available orgs/datasets
- [ ] Client uploads to org server

- [ ] Sentence based mapping and analytics

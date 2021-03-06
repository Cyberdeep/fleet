API Documentation
=================

Kolide Fleet is powered by a Go API server which serves three types of endpoints:

- Endpoints starting with `/api/v1/osquery/` are osquery TLS server API endpoints. All of these endpoints are used for talking to osqueryd agents and that's it.
- Endpoints starting with `/api/v1/kolide/` are endpoints to interact with the Fleet data model (packs, queries, scheduled queries, labels, hosts, etc) as well as application endpoints (configuring settings, logging in, session management, etc).
- All other endpoints are served the React single page application bundle. The React app uses React Router to determine whether or not the URI is a valid route and what to do.

Only osquery agents should interact with the osquery API, but we'd like to support the eventual use of the Fleet API extensively. The API is not very well documented at all right now, but we have plans to:

- Generate and publish detailed documentation via a tool built using [test2doc](https://github.com/adams-sarah/test2doc) (or something competitive).
- Release a JavaScript Fleet API client library (which would be derived from the [current](https://github.com/kolide/fleet/blob/master/frontend/kolide/index.js) JavaScript API client).
- Commit to a stable, standardized API format that we can commit to supporting.

## Current API

The general idea with the current API is that there are many entities throughout the Fleet application, such as:

- Queries
- Packs
- Labels
- Hosts

Each set of objects follows a similar REST access pattern.

- You can `GET /api/v1/kolide/packs` to get all packs
- You can `GET /api/v1/kolide/packs/1` to get a specific pack.
- You can `DELETE /api/v1/kolide/packs/1` to delete a specific pack.
- You can `POST /api/v1/kolide/packs` (with a valid body) to create a new pack.
- You can `PATCH /api/v1/kolide/packs/1` (with a valid body) to modify a specific pack.

Queries, packs, scheduled queries, labels, invites, users, sessions all behave this way. Some objects, like invites, have additional HTTP methods for additional functionality. Some objects, such as scheduled queries, are merely a relationship between two other objects (in this case, a query and a pack) with some details attached.

All of these objects are put together and distributed to the appropriate osquery agents at the appropriate time. At this time, the best source of truth for the API is the [HTTP handler file](https://github.com/kolide/fleet/blob/master/server/service/handler.go) in the Go application. The REST API is exposed via a transport layer on top of an RPC service which is implemented using a micro-service library called [Go Kit](https://github.com/go-kit/kit). If using the Kolide API is important to you right now, being familiar with Go Kit would definitely be helpful.

Like it was said above, we have plans to include richer API documentation in the near future, so stay tuned. If using this API is important to you, please contact us at [support@kolide.co](mailto:support@kolide.co) and tell us, so that we can prioritize creating stable API documentation.

### File Integrity Monitoring

[File Integrity Monitoring](https://osquery.readthedocs.io/en/stable/deployment/file-integrity-monitoring/) can be configured using cURL as illustrated in
the following example. The user must first log in to get an authorization token. This token
must be supplied in the authorization headers of subsequent requests to view or change the FIM configuration.
```shell
# Step 1 Log in
curl -X "POST" "https://localhost:8080/api/v1/kolide/login" \
     -H "Content-Type: text/plain; charset=utf-8" \
     -d $'{
"username": "admin",
"password": "supersecret"
}'
## Login Response
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 17 Aug 2017 22:38:33 GMT
Content-Length: 604
Connection: close

{
  "user": {
    "created_at": "2017-08-17T22:00:32Z",
    "updated_at": "2017-08-17T22:00:32Z",
    "deleted_at": null,
    "deleted": false,
    "id": 1,
    "username": "admin",
    "name": "",
    "email": "admin@acme.com",
    "admin": true,
    "enabled": true,
    "force_password_reset": false,
    "gravatar_url": "",
    "sso_enabled": false
  },
  "token": "faketoken"
}

# Step 2 Upload FIM configuration
curl -X "PATCH" "https://localhost:8080/api/v1/kolide/fim" \
     -H "Authorization: Bearer faketoken" \
     -H "Content-Type: text/plain; charset=utf-8" \
     -d $'{
  "interval": 500,
  "file_paths": {
    "etc": [
      "/etc/%%"
    ],
    "users": [
      "/Users/%/Library/%%",
      "/Users/%/Documents/%%"
    ],
    "usr": [
      "/usr/bin/%%"
    ]
  }
}'

## Upload FIM Response
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 17 Aug 2017 22:39:26 GMT
Content-Length: 3
Connection: close

{}

## View current FIM
curl "https://localhost:8080/api/v1/kolide/fim" \
     -H "Authorization: Bearer faketoken"
```
The FIM configuration supplied in each PATCH request replaces the existing FIM configuration. In order to completely
disable FIM send a PATCH request with an empty set of file paths.
```shell
curl -X "PATCH" "https://localhost:8080/api/v1/kolide/fim" \
     -H "Authorization: Bearer faketoken" \
     -H "Content-Type: text/plain; charset=utf-8" \
     -d $'{
  "interval": 500
}'
```

### Osquery Configuration Import

You can load packs, queries and other settings from an existing [Osquery configuration file](https://osquery.readthedocs.io/en/stable/deployment/configuration/) by importing the file into Fleet. This can be done posting the stringified contents of the Osquery configuration to the following Fleet endpoint:

```
// POST body the value of "config" is JSON that has been converted to a string

{
  "config": "{\"options\":null,\"schedule\":null,\"packs\":{ ...
}

// POST endpoint

/api/v1/kolide/osquery/config/import
```

We provide [a utility program](https://github.com/kolide/configimporter) that will import the configuration automatically.
If you opt to manually import your Osquery configuration you will need to include the contents of externally
referenced packs in your main Osquery configuration file before posting it to Fleet. If you reference packs
in a file like the example below, you will need to get the pack from `external_pack.conf`
and include it in the main configuration.

```
// Configuration referencing external pack

{
  "packs": {
    "external_pack": "/path/to/external_pack.conf",
    "internal_stuff": {
      [...]
    }
  }
}
```

```
// Edited configuration containing the internal pack

{
  "packs": {
    "external_pack": {
      "shard": "10",
      "queries": {
        "suid_bins": {
          "query": "select * from suid_bins;",
          "interval": "3600"
        }
      }
    }
    "internal_stuff": {
      [...]
    }
  }
}
```

Once the configuration file and all the external packs it references are consolidated, post the stringified contents of the configuration file to Fleet.

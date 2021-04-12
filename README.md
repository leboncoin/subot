# Support analytics

Support analytics (aka Subot) is a project designed by leboncoin which aims at simplifying and automating
the support tasks using slack.

## Usage

Users shall ask questions to the support team on a dedicated Slack channel, using threads to organize
discussions by issue. Members of the support team shall respond to the users to acknowledge their issues
and can mark their issues as solved using an emoji on the original message.

The support can have a rotating fireman role. Every week on month one person is in charge of the support.

## Features
- Analytics (analyse common requests, status of the requests, team's velocity, etc...)
- Automatic thread responses (can be basic or based on the content of the message)
- Reminders (recall the fireman after one hour of inactivity on a thread)
- Feedbacks on automatic responses (can lead to automatic solving)
- Reports (send a public report at the end of each week containing the performances of the support team)
- Welcome messages (send ephemeral messages to new members of the channel)

## Architecture

The application has multiple components

### Backend services

- analytics
  Service used by the frontend

- replier
  Service which receives slack webhooks. It also triggers the reporting once a week.

### Storage

- Elasticsearch

### Frontend

See the [frontend project](https://github.com/leboncoin/subot-front)

### Analytics Engine (beta)

See the [analytics engine project](https://github.com/leboncoin/subot-engine)

## Hosting

Multiple hosting methods are possible, but we recommend using a containerized solution like AWS ECS or
a Kubernetes cluster.

## Configuration

Here is the list of all the parameters supported by the application

| Name                              | Environment variable              | Required | Description                                                                                                                                     | Valid values                 | Default                                             |
|-----------------------------------|-----------------------------------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------|-----------------------------------------------------|
| front_url                         | FRONT_URL                         | true     | The URL of the frontend. Used to whitelist for cors                                                                                             |                              |                                                     |
| elasticsearch_url                 | ELASTICSEARCH_URL                 | true     | The URL of the elasticsearch instance.  Elasticsearch shall be up and running prior to running the app                                          |                              |                                                     |
| engine_url                        | ENGINE_URL                        | true     | The URL of the analytics engine which will receive GRPC requests.                                                                               |                              |                                                     |
| analytics_url                     | ANALYTICS_URL                     | true     | The URL at which the analytics service will run.  This is used for the callbacks on the authentication service                                  |                              |                                                     |
| vault_enabled                     | VAULT_ENABLED                     | false    | Boolean to activate vault secret fetching.  Every parameters starting with VAULT::path/to/secret:key  will be read from vault at the given path |                              | false                                               |
| vault_auth_method                 | VAULT_AUTH_METHOD                 | false    | Auth method to use to login into vault if vault is enabled                                                                                      | [token, approle, kubernetes] | token                                               |
| vault_token                       | VAULT_TOKEN                       | false    | Token to use with token auth method.  If not set, will fallback to the users token if available                                                 |                              | defaults to the content of ~/.vault-token           |
| vault_role_id                     | VAULT_ROLE_ID                     | false    | The role for approle login                                                                                                                      |                              |                                                     |
| vault_secret_id                   | VAULT_SECRET_ID                   | false    | The secret id for approle login                                                                                                                 |                              |                                                     |
| vault_approle_mountpoint          | VAULT_APPROLE_MOUNTPOINT          | false    | The path of the login endpoint for the approle auth method                                                                                      |                              | /v1/auth/approle/login                              |
| vault_k8s_token                   | VAULT_K8S_TOKEN                   | false    | The value of the kubernetes token (jwt) to use for kubernetes auth method                                                                       |                              |                                                     |
| vault_k8s_token_path              | VAULT_K8S_TOKEN_PATH              | false    | The path where to look for the kubernetes token (jwt)  to use for kubernetes auth method                                                        |                              | /var/run/secrets/kubernetes.io/serviceaccount/token |
| vault_k8s_role                    | VAULT_K8S_ROLE                    | true     | Name of the role to assume when logging in using kubernetes auth method                                                                         |                              |                                                     |
| vault_k8s_mountpoint              | VAULT_K8S_MOUNTPOINT              | false    | Path of the login endpoint for kubernetes auth method                                                                                           |                              | auth/kubernets/login                                |
| vault_url                         | VAULT_URL                         | false    | The URL of the vault cluster                                                                                                                    |                              | http://localhost:8200                               |
| dex_connectors                    | DEX_CONNECTORS                    | false    | List of connectors for dex authentication                                                                                                       | [ldap]                       | []                                                  |
| dex_admin_group                   | DEX_ADMIN_GROUP                   | true     | Name of the group containing admin members                                                                                                      |                              |                                                     |
| dex_ldap_usersearch_basedn        | DEX_LDAP_USERSEARCH_BASEDN        | true     | Base location of the users allowed to log in                                                                                                    |                              |                                                     |
| dex_ldap_usersearch_filter        | DEX_LDAP_USERSEARCH_FILTER        | true     | Filter for the users - e.g. (objectClass=person)                                                                                                |                              |                                                     |
| dex_ldap_usersearch_username      | DEX_LDAP_USERSEARCH_USERNAME      | true     | Field containing the username of the users (used for login)                                                                                     |                              |                                                     |
| dex_ldap_usersearch_idattr        | DEX_LDAP_USERSEARCH_IDATTR        | true     | Field containing the id of the users                                                                                                            |                              |                                                     |
| dex_ldap_usersearch_emailattr     | DEX_LDAP_USERSEARCH_EMAILATTR     | true     | Field containing the email of the users                                                                                                         |                              |                                                     |
| dex_ldap_usersearch_nameattr      | DEX_LDAP_USERSEARCH_NAMEATTR      | true     | Field containing the name of the users                                                                                                          |                              |                                                     |
| dex_ldap_groupsearch_basedn       | DEX_LDAP_GROUPSEARCH_BASEDN       | true     | Base location for the groups search                                                                                                             |                              |                                                     |
| dex_ldap_groupsearch_filter       | DEX_LDAP_GROUPSEARCH_FILTER       | true     | Filter to apply to the groups search                                                                                                            |                              |                                                     |
| dex_ldap_groupsearch_username     | DEX_LDAP_GROUPSEARCH_USERNAME     | true     | Field for user / group association                                                                                                              |                              |                                                     |
| dex_ldap_groupsearch_emailattr    | DEX_LDAP_GROUPSEARCH_EMAILATTR    | true     | Field containing the email of the groups                                                                                                        |                              |                                                     |
| dex_ldap_groupsearch_nameattr     | DEX_LDAP_GROUPSEARCH_NAMEATTR     | true     | Field containing the name of the group                                                                                                          |                              |                                                     |
| dex_ldap_username                 | DEX_LDAP_USERNAME                 | true     | Username of the user that will be used to browse the ldap                                                                                       |                              |                                                     |
| dex_ldap_password                 | DEX_LDAP_PASSWORD                 | true     | Password of the user that will be used to browse the ldap                                                                                       |                              |                                                     |
| dex_ldap_host                     | DEX_LDAP_HOST                     | true     | URL of the ldap server                                                                                                                          |                              |                                                     |
| dex_secret                        | DEX_SECRET                        | true     | Secret to access the dex server.  We recommend to generate a 16-character random string                                                         |                              |                                                     |
| dex_client_id                     | DEX_CLIENT_ID                     | true     | Unique id for this client.                                                                                                                      |                              |                                                     |
| dex_private_key                   | DEX_PRIVATE_KEY                   | true     | Private key generated to secure communications with the dex server                                                                              |                              |                                                     |
| slack_id                          | SLACK_ID                          | true     | ID of the slack channel to connect the bot to                                                                                                   |                              |                                                     |
| slack_webhook                     | SLACK_WEBHOOK                     | true     | URL to send webhook messages                                                                                                                    |                              |                                                     |
| slack_oauth_access_token          | SLACK_OAUTH_ACCESS_TOKEN          | true     | Oauth access token to the slack API                                                                                                             |                              |                                                     |
| slack_bot_user_oauth_access_token | SLACK_BOT_USER_OAUTH_ACCESS_TOKEN | true     | Oauth access token for the bot user to the API                                                                                                  |                              |                                                     |
| slack_bot_id                      | SLACK_BOT_ID                      | true     | ID of the bot user                                                                                                                              |                              |                                                     |

## Local development

### Start the services

```bash
make start
```

### Stop the services

```bash
make stop
```

### Init the elasticsearch indexes

When starting the elasticsearch container for the first time you will need to initialize its indexes.
To do so, navigate to the infrastructure repository (proctool/infra/support-analytics) and run :
```
cd scripts
./elasticsearch.sh http://localhost:9200 init
```

### Check everything is working

- [Api](http://localhost:8080/ping)
- [Replier](http://localhost:8081/ping)

### Debug in local

To see all the logs, use:

```bash
make logs
```

To see the logs of a specific service, use one of the following commands:

```bash
make log-analytics
make log-replier
```

In case you want to rebuild a service, you can use the following commands:

```bash
make rebuild-api
make rebuild-replier
```

> You may have trouble starting the elasticsearch container.
>
> In case of a `max virtual memory too low`, use the following command to increase it:
>
> ```bash
> sudo sysctl -w vm.max_map_count=262144
> ```

## How to test the webhook

### Send a slack payload

Example of an event from a user joining the channel

```http
POST /event HTTP/1.1
Host: localhost:8081
Content-Type: text/plain

{"token":"<SLACK_TOKEN>","team_id":"<TEAM_ID>","api_app_id":"AL14JDQDQ","event":{"type":"message","subtype":"channel_join","ts":"1577110180.007100","user":"UCMD2JME2","text":"<@UCMD2JME2> has joined the channel","channel":"CLK7MCUS3","event_ts":"1577110180.007100","channel_type":"channel"},"type":"event_callback","event_id":"EvS2E28FFY","event_time":1577110180,"authed_users":["ULGP77XEY"]}
```

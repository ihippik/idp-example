# Identity Provider Example
Toy Identity Provider fo ORY Hydra.

Example for an article on Medium: https://medium.com/scum-gazeta/golang-oauth2-openid-d69d09cb84db

This example shows how to implement a flow - authorization code via
 OAuth 2.0 and OpenID Connect Provider - ORY Hydra.
 
 ### Quick start
 * Run ORY Hydra according to its documentation.
 > docker-compose up
 ```yaml
version: '3'

services:

  hydra:
    image: oryd/hydra:1.4.8
    ports:
      - "4444:4444" # Public port
      - "4445:4445" # Admin port
    command:
      serve all --dangerous-force-http
    environment:
      - URLS_SELF_ISSUER=http://127.0.0.1:4444
      - URLS_CONSENT=http://127.0.0.1:3000/consent
      - URLS_LOGIN=http://127.0.0.1:3000/login
      - URLS_LOGOUT=http://127.0.0.1:3000/logout
      - DSN=memory
      - SECRETS_SYSTEM=youReallyNeedToChangeThis
      - OIDC_SUBJECT_IDENTIFIERS_SUPPORTED_TYPES=public,pairwise
      - OIDC_SUBJECT_IDENTIFIERS_PAIRWISE_SALT=youReallyNeedToChangeThis
    restart: unless-stopped
```
* Create a client that is capable of performing  grant access.
> run in a container with hydra
```shell script
hydra clients create \
    --endpoint http://127.0.0.1:4445 \
    --id scum-client \
    --secret secret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,offline \
    --callbacks http://127.0.0.1:5555/callback
```

* SignIn
> http://127.0.0.1:4444/oauth2/auth?audience=&client_id=scum-client&redirect_uri=http://127.0.0.1:5555/callback&response_type=code&scope=openid+offline&state=pnaqqipwwpbrdkosbqflsnya

As a result, you get an atorization code, which you can then borrow from Hydra for tokens 🎉

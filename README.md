# Minimal REST API

## Project template to fork and implement real rest api

Dependences

- [**gorilla/mux**](https://github.com/gorilla/mux) as router
- [**zerolog**](github.com/rs/zerolog/log) as logger
- Postgresql as db

Branches

- _master_: authentication with JWT autogenerated at startup (see log)
- _authboss-implementation_: authentication with Authboss _**failed...**_, can't route login request

TODO

- authenticaion by provider as Google, Github, ecc
- sessions to manage authorization/expiration
- user info persistence
- CRUD to persist basic entity

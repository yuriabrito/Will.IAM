# Will.IAM

Will.IAM solves identity and access management.

Desired features:

* [ ] Authentication with Google as OAUTH-2 provider.
* [ ] RBAC authorization
  Permissions+Roles+/am
* [ ] SSO - Single Sign-On
  * [ ] SSO browser handler should save/get to/from localStorage and redirect to requester
  Client redirects to server (browser), server has token in localStorage, redirects back with stored token. No button clicks :) Client should be careful to not log token to other parties (e.g google analytics)

## About RBAC use cases and implementation

Client projects of Will.IAM define permissions necessary for resource operation.

Using Maestro, https://github.com/topfreegames/maestro, as an example:

In order to get a list of schedulers, users must have ListSchedulers permission.

Permissions are written in a specific format **{Service}::{OwnershipLevel}::{Action}::{Resource::Hierarchy}**. So, ListSchedulers could be had in a diversity of ways:

Maestro::RO::ListSchedulers::*

Maestro::RL::ListSchedulers::NA::Sniper3D::*

Maestro::RL::ListSchedulers::NA::Sniper3D::sniper3d-game

You'll know more about Will.IAM permissions later. If someone has **Maestro::RL::ListSchedulers::NA::Sniper3D::\***, then Maestro will only respond schedulers under NA::Sniper3D's domain.

## Permissions

Every permission has four components:

### Service

A naming reference to any application service account that uses Will.IAM as IAM solution.

### Ownership Level

**ResourceOwner**: Can exercise the action over the resource and provide the exact
same rights to other parties.

**ResourceLender**: Can only exercise the action over the resource.

### Action

A verb defined by Will.IAM clients.

### Resource Hierarchy

Can be complete or open, in the sense that an open hierarchy will probably lead to access to multiple items under a domain.


## Client side - /am route

Will.IAM clients should expose a **GET /am** route that will help list actions and resource hierarchies to which the requester has some level os access.

E.g:

**GET /am** -> will respond all verbs (actions) the requester has access

**GET /am?permission=ListSchedulers** -> all regions that requester can ListSchedulers

**GET /am?permission=ListSchedulers::NA** -> all games that requester can ListSchedulers in NA

**GET /am?permission=ListSchedulers::NA::Sniper3D** -> all schedulers in NA::Sniper3D

To a requester with full access over the client, this means it will list all possible permissions and resources possible to be granted OwnershipLevel::Action to another party.

## Permission dependency

A nice-to-have feature would be to declare permission dependencies. It should be expected that **Maestro::RL::EditScheduler::\*** implies following **Maestro::RL::ReadScheduler::\***

One way to do this is to have clients declare them over a Will.IAM endpoint and use this custom entity, PermissionDependency, when creating / deleting user|role permissions.

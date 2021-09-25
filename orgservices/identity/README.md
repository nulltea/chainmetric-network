# Chainmetric: Identity Service

User management in a permissioned blockchain infrastructure is complicated and security-intensive enough to be left on user to deal with on its own.

Current service integrates Fabric CA with HashiCorp Vault and exposes gRPC methods for hiding complexity of issuing identity X509 certificates and private key pairs, intricate access control, and blockchain authentication behind plain userpass login and registration scheme.

## API Reference

### [Access service](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/rpc/access_grpc.proto)
#### rpc RequestFabricCredentials([FabricCredentialsRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/access.proto#L10)) returns ([FabricCredentialsResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/access.proto#L15))

Authenticates user with userpass credentials and reserves certificate and private key pair in Vault for accessing Fabric blockchain network.

Generates and returns JWT token for further interacting with API.

#### rpc AuthWithSigningIdentity([CertificateAuthRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/access.proto#L30)) returns ([CertificateAuthResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/access.proto#L35))

Allows authenticate for users that already have Fabric credentials but require JWT token for interacting with current or other org-services.

### [Admin service](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/rpc/admin_grpc.proto)
#### rpc GetCandidates([UsersRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/user.proto#L24)) returns ([UsersResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/user.proto#L26))

Returns all pending-enrollment users of current organization for Admin to review.

#### rpc EnrollUser([EnrollUserRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/admin.proto#L10)) returns ([EnrollUserResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/admin.proto#L16))

Allows Admins to confirm users and allow them to access Fabric blockchain network.

### [User service](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/rpc/user_grpc.proto)

#### rpc Register([RegistrationRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/user.proto#L33)) returns ([RegistrationResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/user.proto#L39))
Allows to submit user registration form, new user identities created for accessing both Fabric network and org-services are initially not active and pending confiramtion from Admin.

#### rpc GetState() returns ([User](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/user.proto#L10))

Returns user data stored on server.

#### rpc ChangePassword([ChangePasswordRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/user.proto#L46)) returns ([StatusResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/shared/proto/status.proto))

Allows to change password.
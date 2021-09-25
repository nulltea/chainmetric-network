# Chainmetric: Identity Service

User management in a permissioned blockchain infrastructure is complicated and security-intensive enough to be left to user to deal with on its own.

Current service integrates Fabric CA with HashiCorp Vault and exposes gRPC methods for hiding complexity of issuing identity X509 certificates and private key pairs, intricate access control, and blockchain authentication behind plain userpass login and registration scheme.

## API Reference

### Access service
#### rpc [RequestFabricCredentials](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/rpc/access_grpc.proto#L19)
Argument: [FabricCredentialsRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/access.proto#L10)

Returns: [FabricCredentialsResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/access.proto#L15)

#### rpc [AuthWithSigningIdentity](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/rpc/access_grpc.proto#L10)
Argument: [CertificateAuthRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/access.proto#L30)

Returns: [CertificateAuthResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/access.proto#L35)

### Admin service
#### rpc [GetCandidates](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/rpc/admin_grpc.proto#L10)
Argument: [UsersRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/user.proto#L24)

Returns: [UsersResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/presenter/user.proto#L26)

#### rpc [EnrollUser](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/identity/api/rpc/admin_grpc.proto#L11)
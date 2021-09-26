# Chainmetric: Notifications Service

For businesses dealing with environmentally-sensitive assets ability to react quickly on incidents and requirements violations is crucial.

Hyperledger Fabric out of the box provide elegant tool set for building an event-driven infrastructure. However, ones established subscription on specific events to continue receiving them, one is required to maintain constant connection to the server.

Such constraint is unacceptable for mobile applications, which for handling push notifications rely on native messaging services like Google's [FCM](https://firebase.google.com/docs/cloud-messaging) or Apple's [APN](https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/APNSOverview.html).

**Notifications service** integrates Firebase Cloud Messaging platform with events sourced from Fabric Orderer services, allows user to subscribe to specific events, and receive push notifications even if their mobile application isn't active.

## Flowchart

![flowchart]

[flowchart]: https://github.com/timoth-y/chainmetric-network/blob/github/update_readme/docs/notifications-service-flowchart.png?raw=true


## API Reference

### [Subscriber service](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/notifications/api/rpc/subscriber_grpc.proto)

#### rpc Subscribe([SubscriptionRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/notifications/api/presenter/subscription.proto#L10)) returns ([SubscriptionResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/notifications/api/presenter/subscription.proto#L26))

Allows users to subscribe to various events on Blockchain. This triggers server to spawn listener routine, on receiving events one would forward them via Firebase API.

#### rpc Cancel([CancellationRequest](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/notifications/api/presenter/subscription.proto#L30)) returns ([StatusResponse](https://github.com/timoth-y/chainmetric-network/blob/main/orgservices/shared/proto/status.proto))

Unsubscribes user from topic related to previously subscribed resource. In case user was the last subscriber of such resource, listener routine is killed too.

## License

Licensed under the [Apache 2.0](https://github.com/timoth-y/chainmetric-network/blob/main/LICENSE).

# meta

## English

Meta is a distributed message control and management framework. Why is it called meta? It means unitary in mathematics. The purpose is to explain that the core design of meta is self-contained, which means that meta can be managed and controlled by other meta, so as to realize infinite data processing and expand downward like a tree. The design concept of meta is extremely advanced.

Note that when we discuss dealing with infinite data, meta does not focus on how to decompose tasks, which should be organized and concerned by the caller, but meta only provides such capabilities.

Meta includes but is not limited to the following functions:
* Manage and control hosts.
* Message type processing.

Architecture: the message distributor distributes the same message to the corresponding processor.
```
            Distribute message type
server -----------------------------------> processor
```

Manager: manager processor. Manage and connect with other servers. You can implement the semantics of for all.

Users can implement their own message types and processors. Processor is a multi process architecture model, which is not the external interface of nodes, but distributed by server, although the processor and server are also communicated by grpc. The purpose of this is to achieve modular decoupling.

Other designs:
* Expansion: it is nothing more than the process of removing and re placing tubes.
* Deployment: the system only focuses on how to manage and deploy easily. How to deploy is decoupled.
* Upgrade: usually, upgrade is related to the deployment mode. For container deployment, there is obviously no way to self upgrade. Because the upgrade should also be decoupled, the only problem to consider is how not to affect the existing system. The entity can enter the maintenance mode when upgrading.

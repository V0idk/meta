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

## Chinese

meta是一个分布式消息控制管理框架。为什么叫meta，元，在数学上有幺元的意思。旨在说明meta的核心设计是自包含的，这意味着meta可以由其他meta管理控制，从而实现无限级的数据处理，就像一颗树一样不断向下扩张。meta的设计理念时是极其先进的。

注意，当我们讨论处理无限的数据时，meta并不关注如何分解任务，这是调用者应该组织和关心的，而meta只是提供这样的能力。

meta包含且不限于以下功能：
* 管理和控制主机。
* 消息类型处理。

架构： 由消息分发器同一分发消息到对应的processor。

```
        根据消息类型分发
server ---------------> processor
```

manager: 管理者处理器。负责管理对接其他server。可以实现for all的语义。

用户可以实现自己的消息类型和processor。processor是多进程架构模型，不是节点对外的接口，而是有server分发，尽管processor和server之间同样是由grpc通信的。这么做的目的是实现模块化的解耦。

其他设计：
* 伸缩: 无非是去纳管和重新纳管的过程。
* 部署: 系统只关注如何管理和容易部署。而如何部署是解耦的。
* 升级：通常升级与部署模式有关。对于容器部署显然没办法做到自升级。因为升级也应该是解耦的，唯一要考虑的问题是如何不影响现有的系统。实体在升级时进入维护模式即可。
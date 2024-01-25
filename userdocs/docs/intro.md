---
title: Introduction
---

# Introduction

:::warning site under construction
:::

## What is flintlock?

Flintlock is a service for creating and managing the lifecycle of microVMs on a
host machine. Wel support [Firecracker][firecracker] and [Cloud Hypervisor][ch].

The primary use case for flintlock is to create microVMs on a bare-metal host
where the microVMs will be used as nodes in a virtualized Kubernetes cluster.
It is an essential part of **Liquid Metal** and can be
driven by [Cluster API Provider Microvm][capmvm].

## Features

Using API requests (via [gRPC][proto] or <a href="/flintlock-api" target="_blank">HTTP</a>):

- Create, update, delete microVMs using Firecracker
- Manage the lifecycle of microVMs (i.e. start, stop, pause)
- Configure microVM metadata via cloud-init, ignition etc
- Use OCI images for microVM volumes, kernel and initrd
- (coming soon) Use CNI to configure the network for the microVMs

## Liquid Metal

To learn more about using Flintlock MicroVMs in a Kubernetes cluster, check
out the [official Liquid Metal docs][lm].


[ch]: https://www.cloudhypervisor.org/
[capmvm]: https://github.com/liquidmetal-dev/cluster-api-provider-microvm
[proto]: https://buf.build/liquidmetal-dev/flintlock
[lm]: https://liquidmetal-dev.github.io/site/
[firecracker]: https://firecracker-microvm.github.io/

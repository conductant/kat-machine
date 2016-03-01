# kat-machine

Simple proof-of-concept for multi-cloud machine provisioning and machine database.

## Why
Machine provisioning and tracking is a pain in a large-cluster, multi-cloud world.
Currently tools like Docker Machine is meant to be used as a desktop tool for machine provisioning and
joining hosts to a Swarm cluster.  The desktop / client side focus of the Docker Machine has some limitations:

+ Operator workflow for provisioning machines and joining machines to the swarm cluster is outside of UCP
+ Docker Cloud uses its own tooling for machine provisioning and bootstrapping and cannot take advantage of the
community support in creating drivers that Docker Machine enjoys.
+ Desktop client stores lots of machine state (e.g. certs, instance ids, data from api calls), and this inventory
information is not accessible by anyone else.
+ At large scale, tracking machine instances, especially across cloud providers, become difficult.

A service that has the broad support of provider drivers as Docker Machine is necessary to be a part of a comprehensive
suite of tools for data center management (virtual, on-prem, or hybrid).

## Requirements
A simple RESTful server that integrates with Docker Machine:

+ To expose endpoints for managing the lifecycles of machines across different cloud providers.
+ To provide endpoints for accessing inventory of machines across providers.
+ To take advantage of all the drivers available in Docker Machine
+ Make all Machine command-line options available via JSON payload.
+ Do so without requiring forks to the Docker Machine code.

## Current State

+ Ability to create, stop, start, restart, and remove machine instances via REST API calls:
  + Support for all drivers supported by Docker Machine -- include `driverName` in the URL
  + Stores driver state in filesystem
+ Support token-based auth so that key endpoints such as machine termination or stop are access controlled.  
  + Server uses signed tokens in API calls.
  + Server depends on another entity to create and sign the auth token.

## TO-DO

+ Add FUSE support so that Machine driver uses familiar file system operations while data is stored in a common 
persistent store (e.g. S3 or etcd)
+ Add API to list all hosts across all providers (drivers).
+ Add indexing by provider, by tags, or any attributes specific to the provider.


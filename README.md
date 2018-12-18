# Storage

Storage is an abstract, distributed, synchronized database orchestrator. 
<br>

`What this means`: Storage will allow you to abstract your data into an `object` format that can then be distributed across any node that completely fulfills the interface.
Using Storage you can then synchronize this data to each node connected to the instance and allow it to orchestrate everything in-between.
<br>

Using Storage is simple, and most of the orchestration will happen behind the scenes depending on the setup. A few example setups ... will be given later ... TODO:

<br>

## Folder structure:
- _impl_: Impl (implementation) contains all of the current implementation that have been approached by Storage.
<br>
<br>

## Current implementations: (ranked in order of priority)
  - `Datastore`: Probably the most mature
  - `Redis`: Probably the most worked on
  - `Multi`: An implementation of storage using multiple backing nodes
  - `GRPC`: An interface for a GRPC server to act as a Storage node using GRPC requests and streams
  - `HTTP`: An interface for a HTTP server to act as a Storage node using requests and websockets
  - `In memory (inmem)`: Basics are done; nothing more than a POC
  - `SQL`: VERY early and VERY low priority
  - `Mongo`: VERY early and VERY low priority
  - `Git`: Later
  - `R-Tree`: Later
  - `GCS`: Later
  - `Couch`: Way later
  - `Cockroach`: Way later
  - `Cassandra`: Way later

<br>

## Plans:
1. As of now there are no metrics collected on each store. In the future I would like to collect metrics and profile the stores to provide better data depending on the situation and the DBs given.

2. In it's current immplementation, Storage has no Raft implementation either so the read repair is manual and is without consensus; this will be done in the future and is a high priority.

3. In the future there will be a separate repo for the `grpc-gateway` interface (`HTTP`/`GRPC` servers) so that the lib is not clutered.

4. Later on, changelogs may be done away with, however, there are some tradeoffs made with every dis/advantage so I need to think about it more. I think the changelogs work much better than an append-only-database as of now, but some thoughts have prevailed in argument against them.
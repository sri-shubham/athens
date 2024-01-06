# Models

Contains all models that will be used for entities used and their storage

- All storage operations are behind interface making it easy to swap databases/implementations if need be in future
- All storage mechanisms push update events to a queue which is in memory queue running with channels but it is also a interface so we can swap it out with any message broker if need be.
- All database entity are using generic utility function to save time
- Everything outside this package will use interfaces to access storage engine

interface <-> rpc <-----> rpc <-> interface

instead of changelogs, just use append only like object store but keep a unique table list of the items
doing this will allow us to keep the speed advantage that using changelogs has but also take advantage of the simplified
architecture that objectstore has. Upserting items will also allow us to go back and recover an item if needed



On Secondary Keys:

You will give the storage interface a map of your secondary values
The db however, only needs to store an array of strings
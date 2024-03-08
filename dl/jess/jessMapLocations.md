A is SOUTH of B and NORTH of C
--------------------
the behavior here has enough complexity to merit some documentation:

1. build the locations as "" kind

2. in GenerateDefaultKinds: try as a "room" -- if it succeeded or duped; its now a room; 
otherwise, assume or ensure that its a door.

3. in GenerateValues:  determine whether the lhs(A) is a room or a door.

## LHS Doors

when A is a door, then, for every rhs element (B):
    
* A door, B door: fail. 
  * both can't be doors.
  
* A door, B room:     
  1. put the door A into room B;
  2. set the compass of B to A ( `B.compass[direction] = A`. )
    ( a room's compass points to a door; there is no destination for the door )

## LHS Rooms

when A is a room, then, for every rhs element (B):

### A room, B door:
ie. room A is direction through door B. `O(B)->A`; and implying `A(p)->O`.
    
if B were a door (ex. in some other room O), we'd want something like:
* `B.destination = A`
* `O.compass[direction] = B`
* `fact: 'dir -> <room>.dir -> B`

we can also manufacture a private door in A that leads into O in the reverse direction; guarding against the case where A already has a door on the reverse side.

to find room O, jess needs to be able to ask about B's whereabouts... after the explicit phrases have been played out. ( GenerateDefaultLocations )

tbd: if this has to be limited to the domain... im think sub-domains write pairs as eval'd changes (rather than initial relations), which means the query wouldn't see them. overall this should probably be limited to the originating domain for simplicity anyway. ( the code would also miss changes to door destinations, etc. in child domains. )

### A room, B room:
room A is direction from room B.

approximately: `B.compass[dir] = A;` and try: `A.compass[reverseDir] = B;`

this can implicitly create two private doors ( guarding against the case where a room already has a door on that particular side; for B to A it's a contradiction if it does. for A to B, that's fine; then do nothing. )

in both cases, also setting the fact of movement from room to room.

reverse directions:
---------------------
the model has an explicit `OppositeOf(dir)`

private door generation:
----------------
generates a door on some particular side of a room; so needs: room & direction.
* set NounKind(door), NounName, scenery, privately named.
* set `room.compass[direction]= door` ( where `name := strings.Replace(room, " ", "-", -1) + "-" + strings.Replace(dir, " ", "-", -1) + "-door"` )

to guard against creating doors where the compass might already be in use, it could first determine the name, and then try to set the compass direction. if that succeeds then actually create the door with that name.

locating:
-------
place the door in the room: `AddNounPair("whereabouts", room, door)`

connecting:
----------
a "connect" function which takes a room, door, desiredDirection, otherRoom.
 adds/checks for conflicting facts;
    probably returns okay if connected , to let caller fail in this case ]
    AddFact("dir", room, desiredDirection, otherRoom)

`room.compass[direction]= door`


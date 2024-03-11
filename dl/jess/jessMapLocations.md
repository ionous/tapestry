R is SOUTH of P and NORTH of C
--------------------
the behavior here has enough complexity to merit some documentation:

1. build the locations as "" kind

2. in GenerateDefaultKinds: try as a "room" -- if it succeeded or duped; its now a room; 
otherwise, assume or ensure that its a door.

3. in GenerateValues:  determine whether the lhs(R) is a room or a door.

## LHS Doors

when R is a door, then, for every rhs element (P):
    
* R door, P door: fail. 
  * both can't be doors.
  
* R door, P room:     
  1. put the door R into room P;
  2. set the compass of P to R ( `P.compass[direction] = R`. )
    ( a room's compass points to a door; there is no destination for the door )

## LHS Rooms

when R is a room, then, for every rhs element (P):

### R room, P door:
ex. room R is direction from door P.

There are two different interpretations:
1. the door is inside the room on its opposite side.
2. the door is in some other room. in that room, moving in the specified direction leads to the door, and the door exits into room R.

Inform only handles the first; but it seems either are valid here.

`O(P)->R`; and implying `R(p)->O`.
    
if P were a door (ex. in some other room O), we'd want something like:
* `P.destination = R`
* `O.compass[direction] = P`
* `fact: 'dir -> <room>.dir -> P`

we can also manufacture a private door in R that leads into O in the reverse direction; guarding against the case where R already has a door on the reverse side.

to find room O, jess needs to be able to ask about P's whereabouts... after the explicit phrases have been played out. ( GenerateDefaultLocations )

tbd: if this has to be limited to the domain... im think sub-domains write pairs as eval'd changes (rather than initial relations), which means the query wouldn't see them. overall this should probably be limited to the originating domain for simplicity anyway. ( the code would also miss changes to door destinations, etc. in child domains. )

### R room, P room:
room R is direction from room P.

approximately: `P.compass[dir] = R;` and try: `R.compass[reverseDir] = P;`

this can implicitly create two private doors ( guarding against the case where a room already has a door on that particular side; for P to R it's a contradiction if it does. for R to P, that's fine; then do nothing. )

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


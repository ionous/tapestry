import createShuttle from "./io.js";
import cmds from "./cmds.js";

// given a valid tapestry command,
// return its signature and body in an array of two elements
function parseCommand(op) {
  for (const k in op) {
    if (k!== "--") {
      return [k, op[k]];
    }
  }
}

// add each item to items recursively
function addItem(allItems, src, parentId) {
  // create if necessary:
  const { id } = src;
  let my = allItems[id];
  if (!my) { 
    my = allItems[id] = src;
  } else {
    // refresh allItems traits
    my.traits.splice(0, Infinity, ...src.traits);
  }
  // update the parent
  my.parentId = parentId;
  // update the children, mapping source data to real data
  my.kids = src.kids.map(srcKid => {
    return addItem(allItems, srcKid, id);    
  });
  return my;
}

let gameOver = false;

// each msg implements the Event slot:
// as of 2023-10-18, the complete set is:
// FrameOutput, PairChanged, SceneEnded, SceneStarted, StateChanged
function processEvent(msg, { allItems, playing }) {
  let out = "";
  const [ sig, args ] = parseCommand(msg);
  switch (sig) {
    // printed text; accumulates over multiple events
    case "FrameOutput:":
      out += args;
      break;

    case "GameSignal:":
      const signal = args;
      console.warn("game signal", signal);
      if (signal === "Quit") {
        gameOver = true;
        playing.value = false;
      }
      break;

    // object state change
    // fix: we need the prev state in order to be able to clear it
    case "StateChanged noun:aspect:prev:trait:":
      const [noun, aspect, prev, trait] = args;
      console.log("state changed", noun, aspect, trait);
      break;

    // relational change
    // fix: we dont get both sides of the relation change:
    // we only get new relations; fine for now.
    case "PairChanged a:b:rel:":
      const [a, b, rel] = args;
      if (rel === "whereabouts") {
        const child = allItems[b];
        // remove from old parent
        const oldParentId = child.parentId;
        if (oldParentId) {
          const oldParent = allItems[oldParentId];
          child.parentId = false;
          if (oldParent.kids) {
            const wasAt = oldParent.kids.findIndex((kid) => kid.id === b);
            if (wasAt >= 0) {
              oldParent.kids.splice(wasAt, 1);
            }
          }
        }
        const newParentId = a;
        if (newParentId) {
          child.parentId = newParentId;
          // if this is a new parent object....
          // we might not have heard about it yet.
          const newParent = allItems[newParentId];
          if (newParent && newParent.kids) {
            newParent.kids.push(child);
          }
        }
        break; // handled whereabouts
      }

    default:
      console.log("unhandled", sig);
  }
  return out;
}


export default class Query {
  constructor({
    shuttle,   // url of api
    narration, // a watched array
    currRoom,
    allItems,
    statusBar, // a watched class Status
    playing
  }) {
    this.statusBar= statusBar;
    this.allItems= allItems;
    this.currRoom = currRoom;
    // fix: rename shuttle?
    this.io = createShuttle(shuttle, (msgs, calls) => {
      let out = "";
      if (typeof msgs === 'string') {
        console.error(msgs);
        return;
      }
      for (let i=0; i< msgs.length; ++i) {
        const msg = msgs[i];
        const call = calls[i];
        //
        const [ sig, body ] = parseCommand(msg);
        switch (sig) {
          // fix: result and events should probably be optional;
          // or, make two commands that satisfy some response interface
          case "Frame result:events:error:":
          {
            const [_res, _evts, error] = body;
            console.warn(error);
            break;
          }
          case "Frame result:events:":
          {
            const [result, events] = body;
            if (events) {
              for (const evt of events) {
                out += processEvent(evt, {
                  allItems,
                  playing,
                });
              }
            }
            if (!call) {
              out += result;
            } else {
              // ick: we debug.Stringify the results to support "any value"
              // so we have to unpack that too.
              const res = result? JSON.parse(result): "";
              call(res);
            }
            break;
          }
          default:
            console.log("unhandled message", sig);
        };
      }
      if (out.length) {
        narration.push(out);
      }
    });
  }
  restart(scene) {
    const { io, statusBar } = this;
    return io.post("restart", scene).then(() => {
      return io.query(
        cmds.storyTitle, (title) => {
          statusBar.title = title;
        },
        cmds.currentScore, (score) => {
          statusBar.score = score;
        },
        cmds.currentTurn, (turn) => {
          statusBar.useScoring = turn >= 0;
          statusBar.turns = turn;
        },
        cmds.locationName, (name) => {
          statusBar.location = name;
        },
        cmds.currentObjects, (root) => {
          this._updateObjects(root);
        }
      );
    });
  }

  fabricate(text) {
    if (gameOver) {
      return Promise.reject(new Error("game over"));
    }
    const player = "self"; // fix: kind of hacky that this is tied to self
    const { io, statusBar, currRoom } = this;
    const prevLoc = currRoom.id;
    return io.query(
      // send player input
      cmds.fabricate(text), null,
      // query for new score and turn each frame
      cmds.currentScore, (score) => {
        statusBar.score = score;
      },
      cmds.currentTurn, (turn) => {
        statusBar.useScoring = turn >= 0;
        statusBar.turns = turn;
      },
      // if there is an event that changes the player's whereabouts
      // re-query location and objects
    ).then(()=> {
      const newLoc = this.allItems[player].parentId;
      return (prevLoc !== newLoc) &&
        io.query(
          cmds.locationName, (name) => {
            statusBar.location = name;
          },
          cmds.currentObjects, (root) => {
            this._updateObjects(root);
          });
    });
  }

  _updateObjects(srcData) {
    const root = addItem(this.allItems, srcData, false);
    this.currRoom.value = root;
  }
}


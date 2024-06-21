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

// each msg implements the Event slot:
// as of 2023-10-18, the complete set is:
// FrameOutput, PairChanged, SceneEnded, SceneStarted, StateChanged
function processEvent(objCatalog, msg) {
  let out = "";
  const [ sig, args ] = parseCommand(msg);
  switch (sig) {
    // printed text; accumulates over multiple events
    case "FrameOutput:":
      out += args;
      break;

    // object state change
    // fix: we need the prev state in order to be able to clear it
    case "StateChanged noun:aspect:trait:":
      const [noun, aspect, trait] = args;
      console.log("state changed", noun, aspect, trait);
      break;

    // relational change
    // fix: we dont get both sides of the relation change:
    // we only get new relations; fine for now.
    case "PairChanged a:b:rel:":
      const [a, b, rel] = args;
      if (rel === "whereabouts") {
        const child = objCatalog.get(b);
        // remove from old parent
        const oldParentId = child.parentId;
        if (oldParentId) {
          const oldParent = objCatalog.get(oldParentId);
          child.parentId = false;
          if (oldParent.contents) {
            const wasAt = oldParent.contents.findIndex((el) => el.id === b);
            if (wasAt >= 0) {
              oldParent.contents.splice(wasAt, 1);
            }
          }
        }
        const newParentId = a;
        if (newParentId) {
          child.parentId = newParentId;
          // if this is a new parent object....
          // we might not have heard about it before --
          // we'll probably hear about it in rebuild
          // alt: we could could create a temp folder item here.
          const newParent = objCatalog.get(newParentId);
          if (newParent && newParent.contents) {
            newParent.contents.push(child);
          }
        }
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
    statusBar, // a watched class Status
    objCatalog // class ObjectCatalog
  }) {
    this.statusBar= statusBar;
    this.objCatalog= objCatalog;
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
                out += processEvent(objCatalog, evt);
              }
            }
            if (call) {
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
    const io = this.io;
    const objCatalog = this.objCatalog;
    const statusBar = this.statusBar;
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
           objCatalog.rebuild(root);
        }
      );
    });
  }

  fabricate(text) {
    const player = "self"; // fix: kind of hacky that this is tied to self
    const io = this.io;
    const statusBar= this.statusBar;
    const objCatalog = this.objCatalog;
    const prevLoc = objCatalog.get(player).parentId;
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
      const newLoc = objCatalog.get(player).parentId;
      return (prevLoc !== newLoc) &&
        io.query(
          cmds.locationName, (name) => {
            statusBar.location = name;
          },
          cmds.currentObjects, (objs) => {
            objCatalog.rebuild(objs);
          });
    });
  }
}


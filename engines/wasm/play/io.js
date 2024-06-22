import http from '/lib/http.js'

// fix? rename Shuttle? or ShuttleJs
// fix: maybe could include different library packages or files?
export default function(url, cb) {
   const io = url === "wasm" ? 
               new Wasm(cb) : 
               new Http(url, cb);
   return {
    // assumes an endpoint and some data for that endpoint
    // ex. the query endpoint takes an array of json commands
    post: io.post.bind(io),
    // assumes an array of alternating commands and callbacks.
    // [n]   = a tapestry command ( not serialized yet )
    // [n+1] = the callback to process that result of that command.
    query(...msgCalls) {
      if (msgCalls & 1) {
        throw new Error("expected an equal number of queries and calls");
      }
      const sends = [];
      const calls = [];
      for (let i=0; i<msgCalls.length; i+=2) {
        const send = msgCalls[i];
        const call = msgCalls[i+1];
        sends.push(send);
        calls.push(call);
      }
      // send an array of send; expects an array back.
      return io.post("query", sends, calls);
    }
  }
}

class Wasm {
  constructor(msgcb) {
    if (!tapestry) {
      throw Error("couldn't find the global tapestry object");
    }
    this.tapestry= tapestry;
    this.msgcb= msgcb;
  }
  // data is some pod-like json data.
  // calls is a matching array of callbacks for those commands
  post(endpoint, data, calls) {
    const send = JSON.stringify(data);
    return this.tapestry.post(endpoint, send).then((frames) => {
      return this.msgcb(frames, calls || []);
    }).catch((e)  => {
      console.warn("io error", e);
    });
  }
};

class Http {
  // pass in the sink we're writing data to
  constructor(url, msgcb) {
    this.url= url;
    this.msgcb= msgcb;
  }
  // send a tapestry command
  // promises to call msgcb with the frame, and yields its result.
  post(endpoint, send, calls) {
    console.log("posting", endpoint, send);
    return http.post(this.url+endpoint, send).then((frames) => {
      return this.msgcb(frames, calls || []);
    }).catch((e)  => {
      console.warn("io error", e);
    });
  }
};

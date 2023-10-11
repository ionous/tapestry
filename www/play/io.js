import http from '/lib/http.js'

export default class Io {
  // pass in the sink we're writing data to
  constructor(url, msgcb) {
    this.url= url;
    this.msgcb= msgcb;
  }
  // send a tapestry command
  post(endpoint, send, calls) {
    console.log("posting", endpoint, send);
    return http.post(this.url+endpoint, send).then((frames) => {
      return this.msgcb(frames, calls || []);
    }).catch((e)  => {
      console.warn("io error", e);
    });
  }
  // frames
  query(msgCalls) {
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
    return this.post("query", sends, calls);
  }
} // Io class

import http from '/lib/http.js'

export default class Io {
  // pass in the sink we're writing data to
  constructor(output, endpoint) {
    this.endpoint= endpoint;
    this.keepalive=0;
    this.timer=0;
    this.lines= output;
    this.getting= false;
    this.sending= false;
  }
  stopPolling() {
    this._clearTimer();
    this.keepalive= false;
  }
  startPolling(keepalive= 15) {
    this.stopPolling();
    if (keepalive > 0) {
      this.keepalive= keepalive;
      this.timer = setTimeout( keepalive*1000, () => {
        this.timer=0; // done, so forget.
        this._poll();
      });
    }
  }
  // send a pod command
  send(cmd) {
    this._clearTimer();
    // we use promises to keep our get/send requests "serialized" --
    // but it shouldnt really matter because the server needs to handle things regardless.
    this.sending= Promise.allSettled([this.getting, this.sending]).then(()=>{
      // expects zero or more messages back
      return http.post(this.endpoint, cmd, true).then((msgs) => {
        if (Array.isArray(msgs)) {
          for (const msg of msgs) {
            this.lines.push("from server: "+ msg);
          }
        }
      }).finally(()=>{
        this.startPolling(this.keepalive);
      });
    });
  }
  // poll resets the timer
  // expects zero or more messages
  _poll() {
    this.getting= Promise.allSettled([this.getting, this.sending]).then(()=>{
      return http.get(this.endpoint).then((msgs) => {
        for (const msg of msgs) {
          this.lines.push(msg);
        }
      }).finally(()=>{
        this.startPolling(this.keepalive);
      });
    });
  }
  _clearTimer() {
    if (this.timer) {
      clearTimeout(this.timer);
      this.timer=0;
    }
  }
}

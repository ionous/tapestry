import http from '/lib/http.js'

export default class Io {
  // pass in the sink we're writing data to
  constructor(endpoint, msgcb) {
    this.endpoint= endpoint;
    this.keepalive=0;
    this.timer=0;
    this.msgcb= msgcb;
    this.getting= false;
    this.sending= false;
  }
  stopPolling() {
    this._clearTimer();
    this.keepalive= -1;
  }
  startPolling(keepalive= 15, quick=false) {
    this.stopPolling();
    if (keepalive >= 0) {
      this.keepalive= keepalive;
      this.timer = setTimeout( () => {
        this.timer=0; // done, so forget.
        this._poll();
      }, quick?0: keepalive*1000);
    }
  }
  // send a pod command
  send(cmd) {
    this._clearTimer();
    // we use promises to keep our get/send requests "serialized" --
    // but it shouldnt really matter because the server needs to handle things regardless.
    let msgCnt=0;
    this.sending= Promise.allSettled([this.getting, this.sending]).then(()=>{
      // expects zero or more messages back
      return http.post(this.endpoint, cmd).then((msgs) => {
        if (Array.isArray(msgs)) {
          this.msgcb(msgs);
          msgCnt= msgs.length;
        }
      }).catch((e)  => {
        console.warn("io error", e);
      }).finally(()=>{
        this.startPolling(this.keepalive, msgCnt>500);
      });
    });
  }
  // poll resets the timer
  // expects zero or more messages
  _poll() {
    let msgCnt=0;
    this.getting= Promise.allSettled([this.getting, this.sending]).then(()=>{
      return http.get(this.endpoint).then((msgs) => {
        this.msgcb(msgs);
        msgCnt= msgs.length;
      }).catch((error) => {
        console.log('error:', error)
      }).finally(()=>{
        this.startPolling(this.keepalive, msgCnt>500);
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

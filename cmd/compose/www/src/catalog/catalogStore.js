class Cataloger {
  constructor(base = "") {
    this.base= base;
  }
  _get(path, cb) {
    this._send("GET", path, cb);
  }
  _put(path, body, cb) {
    this._send("PUT", path, cb, body);
  }
  _post(path, body, cb) {
    this._send("POST", path, cb, body);
  }
  _send(method, path, cb, body) {
    const url= this.base+path;
    console.log("xml http request:", method, url);
    var xhr = new XMLHttpRequest();
    xhr.addEventListener("load", ()=>{
      console.log("xml http response:", method, url, xhr.statusText);
      let data= true;
      if (!xhr.response || xhr.response.startsWith('<')) {
        data= xhr.response;
      } else{
        try {
          data= JSON.parse(xhr.response);
        } catch (e) {
          data= false;
        }
      }
      cb(data);
    });
    xhr.addEventListener("abort", ()=>cb(false));
    xhr.addEventListener("error", ()=>cb(false));
    xhr.open(method, url);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(body);
  }
}


// catalog store tracks individual files.
class CatalogStore {
  constructor(nodes) {
    this.nodes= nodes;
    this.cache= {};
  }
  getStory(path) {
    return this.cache[path];
  }
  storeStory(path, storyData) {
    const { nodes, cache } = this;
    const story= nodes.unroll(storyData);
    cache[path]= story;
    return story;
  }
  toJSON() {
    const { cache } = this;
    return Object.keys(cache).map(path=> {
      const story= cache[path];
      return {
        path,
        story,
      };
    });
  }
}

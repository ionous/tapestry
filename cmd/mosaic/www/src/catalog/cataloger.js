// baseclass for list of files and folders
export default class Cataloger {
  constructor(base = "") {
    this.base= base;
  }
  loadFolder(folderItem) {
    throw new Error("load folder not implemented");
  }
  loadFile(fileItem) {
    throw new Error("load file not implemented");
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


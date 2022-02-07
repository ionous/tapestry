import Cataloger from './cataloger.js'
import { CatalogFolder, CatalogFile } from './catalogItems.js'
import { reactive } from "vue";

export default class BlockCatalog extends Cataloger {
  constructor(url) {
    super(url);
    this._store= {}; // stores catalog items
    this._saving= false;
  }
  get busy() {
    return this._saving;
  }
  saveStories() {
    // if (!this._saving) {
    //   this._saving= true;
    //   const json= JSON.stringify(this._store);
    //   this._put("", json, (res)=>{
    //     console.log("SAVED:", res);
    //     this._saving= false;
    //   });
    // }
  }
  // inject a directory listing into folder
  loadFolder(folder) {
    const { path } = folder;
    const { _store: store } = this;
    this._get(path).then((contents) => {
      if (Array.isArray(contents)) {
        // a list of CatalogItems
        folder.contents= this._readFolder(path, contents);
        for (const el of folder.contents) {
          store[el.path]= el;
        }
      }
    });
  }
  // inject the contents of a file into itself
  // promise the new contents
  loadFile(path) {
    let ret;
    const { _store: store } = this;
    let file= store[path];
    if (file && file.contents) {
      ret= Promise.resolve(file);
    } else {
      // FIX: can we *not* parse the json?
      // or make make it _contents
      ret= this._get(path).then((result) => {
        console.log('success:', path);
        let file= store[path];
        if (!file) {
          const parts= path.split('/');
          const dir= parts.slice(0, parts.length-1).join("/");
          const name= parts[parts.length-1];
          console.log("adding file @", name, dir);
          store[path]= file= new CatalogFile(name, dir);
        }
        file.contents= Object.freeze(result);
        return file;
      });
    }
    return ret;
  }
  // run an action like "check", etc. against a specific file.
  run(action, file, options, cb) {
    const { path } = file;
    this._post(`${path}/${action}`, options, cb);
  }
  // turns a json array of paths into CatalogItems
  // ["/curr","/proj1","/proj2","/shared", "currStory.if"]
  _readFolder(path, got) {
    console.log("gotten:", got);
    return got.map((el)=> {
      const isFolder= el.startsWith("/");
      const name= el.slice(isFolder?1:0);
      let ret= isFolder? new CatalogFolder(name, path): new CatalogFile(name, path);
      return reactive(ret);
    });
  }
  // promise json
  _get(path) {
    const url= this.base + path;
    console.log("getting", url);
    return fetch(url).then((response) => {
      return (!response.ok) ?
        Promise.reject({status: response.status, url}) :
        response.json().then(result => result);
    }).catch((error) => {
      console.log('error:', error)
    });
  }
};

import Cataloger from './cataloger.js'
import http from './http.js'
import { CatalogFolder, CatalogFile } from './catalogItems.js'
import { reactive } from "vue"

export default class BlockCatalog extends Cataloger {
  constructor(url) {
    super(url);
    this._store= {}; // stores catalog items
    this._saving= false; // false or an arbitrary non-zero number
    this._saved= false; // false or an arbitrary non-zero number
    this._saves= 0; //
  }
  get busy() {
    return !!this._saving;
  }
  // inject an .if directory listing into folder
  loadFolder(folder) {
    const { path } = folder;
    const { _store: store } = this;
    http.get(this.base, path).then((contents) => {
      if (Array.isArray(contents)) {
        // a list of CatalogItems
        folder.contents= this._readFolder(path, contents);
        for (const el of folder.contents) {
          if (!store[el.path]) { // fix? what if something gets deleted?
            store[el.path]= el;
          }
        }
      }
    });
  }
  // store some json string content in memory.
  // *doesnt* save it to the server.
  storeFile(path, contents) {
    const file= this._store[path];
    // we expect that to store a file, we already have a file.
    if (!file || !file.contents) {
      throw new Error(`file ${path} never loaded`);
    } else {
      file.updateContents(contents);
    }
  }
  // inject the contents of a file into itself
  // promise the new contents
  loadFile(path) {
    let ret;
    const { _store: store } = this;
    let file= store[path];
    if (file && file.contents) {
      // already exists? return a promise to match
      ret= Promise.resolve(file);
    } else {
      ret= http.get(this.base, path).then((pod) => {
        console.log('success:', path);
        let file= store[path];
        if (!file) {
          const parts= path.split('/');
          const dir= parts.slice(0, parts.length-1).join("/");
          const name= parts[parts.length-1];
          console.log("adding file @", name, dir);
          store[path]= file= reactive(new CatalogFile(name, dir));
        }
        const contents= JSON.stringify(pod);
        file.updateContents(contents, true);
        return file;
      });
    }
    return ret;
  }
  // collect all changes from all files and put them to the server.
  saveStories() {
    const { _saving:saving, _saved:last, _store:store } = this;
    if (!saving) {
      const next= ++this._saves; // increment the save attempt.
      this._saving = next; // stores the save id for debugging.
      //
      let out= [];
      for (const key in store) {
        const item= store[key];
        if (item instanceof CatalogFile) {
          const contents= item.collect(last, next);
          if (contents) {
            out.push({
              path: item.path,
              contents: {toJSON: ()=> contents},
            });
          }
        }
      }
      if (out.length>0) {
        console.log(json);
        http.put(this.base, '', out).then((res)=>{
          console.log("SAVED:", res);
          this._saved= next;
          this._saving= false;
        }).catch((reason)=>{
          console.log("failed to save", reason);
          this._saving= false;
        });
      }
    }
  }
  // run an action like "check", etc. against a specific file.
  // run(action, file, options, cb) {
  //   const { path } = file;
  //   this._post(`${path}/${action}`, options, cb);
  // }

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
};

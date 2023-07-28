import Cataloger from './cataloger.js'
import { CatalogFolder, CatalogFile } from './catalogItems.js'
import http from '/lib/http.js'

// fix? i dont know why these were reactive; i dont think they need to be
//import { reactive } from 'vue'
function reactive(r) {
  return r;
}

export default class BlockCatalog extends Cataloger {
  constructor(url) {
    super(url);
    this._store= {}; // stores catalog items
  }
  // inject an .if directory listing into folder
  loadFolder(folder) {
    const { path } = folder;
    const { _store: store } = this;
    http.get(http.join(this.base, path)).then((contents) => {
      if (Array.isArray(contents)) {
        // a list of CatalogItems
        folder.contents= this._readFolder(path, contents);
        for (const el of folder.contents) {
          if (!store[el.path]) { // fix? what if something gets deleted?
            store[el.path]= el;
          }
        }
      }
    }).catch((error) => {
      console.log('error:', error)
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
      ret= http.get(http.join(this.base, path)).then((pod) => {
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
        file.updateContents(contents);
        return file;
      });
    }
    return ret;
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
};

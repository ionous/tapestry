import { CatalogFolder } from './catalogItems.js'

// baseclass for list of files and folders
export default class Cataloger {
  constructor(base = "") {
    this.base= base;
    this.root= new CatalogFolder("");
  }
  loadFolder(folderItem) {
    throw new Error("load folder not implemented");
  }
  loadFile(path) {
    throw new Error("load file not implemented");
  }
}


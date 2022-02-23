// base-class for all catalog items
class CatalogItem {
  constructor(name, dir) {
    this.name= name;
    this.dir= dir; // ex. curr/sub
    this.contents= false;
  }
  // ex. "", "curr", "curr/sub"
  get path() {
    const { dir, name } = this;
    return dir? `${dir}/${name}`: name;
  }
}

// contents is an array of catalog items
export class CatalogFolder extends CatalogItem {}

// contents is a json string ( to make change detection easier )
export class CatalogFile extends CatalogItem {
  constructor(name, dir) {
    super(name,dir);
    this._changed= false; // can be a bool, or a guid string.
  }
  // returns the contents if there are any pending changes
  // ( as determined by comparing the lastChange marker )
  collect(lastSave, nextSave) {
    let ret= false;
    const{ _changed: changed } = this;
    if (changed) {
      if (changed === lastSave) {
        this._changed= false;
      } else {
        this._changed= nextSave;
        ret= this.contents;
      }
    }
    return ret;
  }
  // changed should be true or false
  // depending on whether the content was just loaded or not.
  updateContents(contents, isSaved=false) {
    if (typeof contents !== 'string') {
      throw new Error("expected a string for file contents");
    }
    if (contents === this.contents) {
      console.log("no change for", this.path);
    } else {
      if (!this.contents) {
        console.log("new contents for", this.path);
      } else {
        console.log("updating contents for", this.path);
      }
      this.contents= contents;
      this._changed= !isSaved;
    }
  }
}

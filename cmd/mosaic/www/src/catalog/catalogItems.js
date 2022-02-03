// base-class for all catalog items
class CatalogItem {
  constructor( name, dir ) {
    this.name= name;
    this.dir= dir; // ex. curr/sub
  }
  // ex. "", "curr", "curr/sub"
  get path() {
    const { dir, name } = this;
    return dir? `${dir}/${name}`: name;
  }
}

export class CatalogFolder extends CatalogItem {
  constructor( name, dir ) {
    super(name, dir);
    this.contents= false;
  }
}

export class CatalogFile extends CatalogItem {
  constructor( name, dir ) {
    super(name, dir);
    this.story= false;
  }
}


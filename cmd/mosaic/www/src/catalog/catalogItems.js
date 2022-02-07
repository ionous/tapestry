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

export class CatalogFolder extends CatalogItem {}
export class CatalogFile extends CatalogItem {}


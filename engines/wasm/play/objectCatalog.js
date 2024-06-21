import Cataloger from '/catalog/cataloger.js'
import { CatalogFolder, CatalogFile } from '/catalog/catalogItems.js'

const mockObjects = {
  id: "tower",
  name: "tower",
  kind: "rooms",
  traits: ["lit"],
  kids: [{
    id: "mirror",
    name: "mirror",
    kind: "things",
    traits: ["fixed in place"],
    kids: []
  },{
    id: "apple",
    name: "apple",
    kind: "containers",
    traits: ["open","portable"],
    kids: [{
      id: "worm",
      name: "worm",
      kind: "actors",
      traits: ["unhappy"],
      kids: []
    }]
  }]
};

// we assume that certain kinds have children
// ( it'd be nice if the tree view hadn't made that decision from the start )
const kindOfLocale = ["rooms", "containers", "supporters", "actors"];

// recursively create an object tree
function buildTree(all, parentId, one) {
  let out;
  if (kindOfLocale.indexOf(one.kind) == -1) {
    const item = new CatalogFile(one.name, parentId);
    item.data = one;
    all[one.id] = item;
    out = item;
  } else {
    const folder = new CatalogFolder(one.name, parentId);
    folder.data = one;
    folder.contents = [];
    all[one.id] = folder;
    out = folder;
    // recurse ( even if its a room,
    if (one.kids) {
      for (let i=0; i< one.kids.length; i++) {
        let child = buildTree(all, one.id, one.kids[i]);
        folder.contents.push(child);
        child.parentId = one.id;
      }
    }
  }
  return out;
}

export default class ObjectCatalog extends Cataloger {
  constructor() {
    super();
    this.all = {};
    // my tree control is ... odd
    // among other behaviors: it hides the root folder
    // and only displays its contents; fine but we actually want that.
    this.root = new CatalogFolder("$root");
    this.root.contents= [];
    this.rebuild(mockObjects);
  }
  get(id) {
    return this.all[id];
  }
  // objs is a json'd "object collection" ( from collect.if )
  rebuild(collection) {
    const root = buildTree(this.all, "", collection);
    this.root.contents.splice(0, 1, root);
    return this.root;
  }
};

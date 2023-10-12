import Cataloger from '/mosaic/catalog/cataloger.js'
import { CatalogFolder, CatalogFile } from '/mosaic/catalog/catalogItems.js'

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

const kindOfLocale = ["containers", "supporters", "actors", "rooms"];

// recursively create an object tree
function addTree(all, parentId, one) {
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
    // recurse:
    for (let i=0; i< one.kids.length; i++) {
      let child = addTree(all, one.id, one.kids[i]);
      folder.contents.push(child);
    }
  }
  return out;
}

export default class ObjectCatalog extends Cataloger {
  constructor() {
    super();
    this.all = {};
    this.root = addTree(this.all, "", mockObjects);
  }

  loadFolder(folder) {
    // TODO: build from all, and/or make some queries
    // if (!folder.contents) {
    //   const contents=[];
    //   for (const k in this.all) {
    //     const item= this.all[k];
    //     if (item.dir == folder.path) {
    //       contents.push(this.all[k]);
    //     }
    //   };
    //   folder.contents= contents;
    // }
  }
};

import Cataloger from './cataloger.js'
import { CatalogFolder, CatalogFile } from './catalogItems.js'

const mockCatalog= [{
  "curr": [
    "currStory.if", {
      "alt": [
        "fileA.if",
        "fileB.if"
      ],
      "sub": [
        "fileC.if",
        "fileD.if",
        "fileE.if"
      ],
    },
  ],
  "proj": [
    "fileA.if",
    "fileB.if",
    "fileC.if",
    "fileD.if",
    "fileE.if"
  ],
  "empty": []
}];

function addAll(parent, all, list, path) {
 for (let i=0; i< list.length; i++) {
  const el= list[i];
  if (typeof el === 'string') {
    const file= new CatalogFile(el, path);
    all[file.path]= file;
   } else {
      for (const name in el) {
        const folder= new CatalogFolder(name, path);
        all[folder.path]= folder;
        const subPath= path? `${path}/${name}`: name;
        addAll(folder, all, el[name], subPath);
      }
    }
  }
  return all;
}

export default class MockCatalog extends Cataloger {
  constructor() {
    super();
    this.store= {};
    this.root= new CatalogFolder("");
    this.all= addAll(this.root, {}, mockCatalog, "");
  }

  // future:takes a string, list, or none; if none saves all.
  saveStories() {
    const json= JSON.stringify(this.store, 0,2);
    console.log("SAVED:", json);
  }

  findByPath(path) {
    return this.all[path];
  }

  loadFile(file) {
    const path= file.path;
    let story= this.store[path];
    if (!story) {
      // previously created an instance of a random story statement
      // const pick= Types.slats("story_statement");
      // const order= path.split("").reduce((a,v)=>(a+v.charCodeAt(0)), 0);
      // // const type= pick[ Math.floor(Math.random() * pick.length) ];
      // const type= pick[order%pick.length];
      // const storyData= Types.createItem(type.name);
      // story= this.store[path]= storyData;
      const storyData= {"blocks": {"languageVersion": 0,"blocks": []}};
      story= this.store[path]= storyData;
    }
    file.contents= story;
  }

  loadFolder(folder) {
    if (!folder.contents) {
      const contents=[];
      for (const k in this.all) {
        const item= this.all[k];
        if (item.dir == folder.path) {
          contents.push(this.all[k]);
        }
      };
      folder.contents= contents;
    }
  }
};

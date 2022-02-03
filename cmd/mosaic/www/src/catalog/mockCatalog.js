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

export default class MockCatalog extends Cataloger {
  constructor() {
    super();
    this.store= {};
  }

  // future:takes a string, list, or none; if none saves all.
  saveStories() {
    const json= JSON.stringify(this.store, 0,2);
    console.log("SAVED:", json);
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
    let mockFolder= mockCatalog;
    const path= folder.path;
    if (path.length) {
      // advance through the mockup data
      path.split("/").forEach((part)=>{
        // the last element is where the directories live.
        const last= mockFolder[mockFolder.length-1];
        const dirs= (typeof last !== 'string')? last: {};
        if (!part in dirs) {
          throw new Error(`unknown part ${part} in ${dirs}`);
        }
        mockFolder= dirs[part];
      });
    }
    folder.contents= MockCatalog.buildContents(mockFolder, path);
    // console.log("folder", folder.name, "has", folder.contents?folder.contents.length:0, "items");
  }

  // returns a list of catalog folders and files.
  static buildContents(mockFolder, mockPath) {
    const ret= [];
    // walk the items in the folder:
    for (const item of mockFolder) {
      // file vs sub-folder map
      if (typeof item === 'string') {
        ret.push( new CatalogFile(item, mockPath) );
      } else {
        // note: if there's a map, its always the last item of the folder
        Object.keys(item).forEach((mockDir) => {
          ret.push( new CatalogFolder(mockDir, mockPath) );
        });
      }
    }
    return ret;
  }
};

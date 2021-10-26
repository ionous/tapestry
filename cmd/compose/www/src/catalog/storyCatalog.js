class StoryCatalog extends Cataloger {
  constructor(nodes) {
    super("/stories/");
    this.store= new CatalogStore(nodes);
    this._saving= false;
    this._pending= false; // a pending CatalogFile
    // fix: really we should have a global statemachine with states like
    // loading, saving, and editing.
  }
  get busy() {
    return this._saving || this._pending;
  }
  saveStories() {
    if (!this._saving) {
      this._saving= true;
      const json= JSON.stringify(this.store);
      this._put("", json, (res)=>{
        console.log("SAVED:", res);
        this._saving= false;
      });
    }
  }
  // inject a directory listing into folder
  loadFolder(folder) {
    const { path } = folder;
    this._get(path, (contents)=>{
      if (Array.isArray(contents)) {
        folder.contents= this.readContents(path, contents);
      }
    });
  }
  // inject a story file into folder.
  loadStory(file) {
    // if we are already loading something
    // set the desired thing to load as the new thing
    // for use when load completes.
    if (this._pending) {
      this._pending= file;
    } else {
      const { path } = file;
      let story= this.store.getStory(path);
      if (story) {
        file.story= story;
      } else {
        this._pending= file;
        this._get(path, (storyData)=>{
          if (storyData) {
            story= this.store.storeStory(path, storyData);
            file.story= story;
            const reload= this._pending;
            this._pending= false;
            this.loadStory(reload);
          }
        });
      }
    }
  }
  run(action, file, options, cb) {
    const { path } = file;
    this._post(`${path}/${action}`, options, cb);
  }
  // ["/curr","/proj1","/proj2","/shared", "currStory.if"]
  readContents(path, got) {
    console.log("gotten:", got);
    return got.map((el)=> {
      const isFolder= el.startsWith("/");
      const name= el.slice(isFolder?1:0);
      return isFolder? new CatalogFolder(name, path): new CatalogFile(name, path);
    });
  }
};

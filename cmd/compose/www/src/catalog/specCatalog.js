class SpecCatalog extends Cataloger {
  constructor(nodes) {
    super("/ifspec/")
    this.specs = {};
    this.maker = new Make(new Types());
  }
  load(allDone) {
    this._get("", (res)=>{
      const regexp = /"(\w+)\.ifspec"/g;
      // each match is an array of [matched text, each group...]
      const matches = Array.from(res.matchAll(regexp)).map(match=>match[1]);
      this._loopingLoad(matches, allDone);
    });
  }
  _loopingLoad(matches, allDone) {
    if (!matches.length) {
      this._allDone(allDone);
    } else {
      const spec= matches.pop();
      this._get(spec + '.ifspec', (res)=>{
        this.specs[spec]= res;
        this._loopingLoad(matches, allDone);
      });
    }
  }
  _allDone(allDone) {
    for (const k in this.specs) {
      const spec= this.specs[k];
      this.maker.readSpec(spec, k);
    }
    console.log(JSON.stringify(this.maker.types.all,0,2));
    if (allDone) {
      allDone();
    }
  }
}

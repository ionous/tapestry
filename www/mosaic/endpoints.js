
export default {
  play: appcfg.mosaic, // post endpoint
  blocks: appcfg.mosaic + '/blocks/',
  // shapes is an array of blockly shape data.
  shapes: appcfg.mosaic + '/shapes/',
  // toolbox data
  tools:  appcfg.mosaic + '/boxes/',

  file(which) {
    return this.blocks + which;
  },
  action(which) {
    return appcfg.mosaic + '/actions/' + which;
  },
};

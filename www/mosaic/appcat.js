import BlockCatalog from './catalog/blockCatalog.js'
import endpoint from './endpoints.js'

// catalog holds the lazy loading list of all possible .if files.
// appcfg comes through vite conifg.
const catalog = new BlockCatalog(endpoint.blocks);

export default catalog;

<template>
<ol class="mk-folder-ctrl"
  ><FolderItem
    v-for="subFolder in folders"
    :key="subFolder.path"
    :folder="subFolder"
    :depth="depth"
    @activated="onFolder(subFolder)"
    ><Folder
      :folder="subFolder"
      :depth="depth+1"
    ></Folder
  ></FolderItem
  ><FileItem
    v-for="file in files"
    :key="file.path"
    :file="file"
    :depth="depth"
    @activated="onFile(file)"
  ></FileItem
></ol>
</template>
<script>

import FileItem from './FileItem.vue'
import FolderItem from './FolderItem.vue'
import { CatalogFolder } from './catalogItems.js'

export default {
  inject: ['onFolder', 'onFile'],
  components: { FolderItem, FileItem },
  props: {
    folder: CatalogFolder,
    depth: {
      type: Number,
      default: 0
    }
  },
  computed: {
    folders() {
      return this.items(true);
    },
    files() {
      return this.items(false);
    }
  },
  methods: {
    items(isFolder) {
      const { folder } = this;
      return folder.contents? folder.contents.filter((el)=> {
        return (el instanceof CatalogFolder) === isFolder;
      }).sort((a,b)=> a.name.localeCompare(b.name)): [];
    }
  }
}
</script>

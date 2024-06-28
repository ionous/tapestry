<template>
<ul class="mk-folder-ctrl"
  ><mk-folder-item
    v-for="subFolder in folders"
    :key="subFolder.path"
    :folder="subFolder"
    :depth="depth"
    @activated="onFolder(subFolder)"
    ><mk-folder
      :folder="subFolder"
      :depth="depth+1"
      @fileSelected="onFile"
      @folderSelected="onFolder"
    ></mk-folder
  ></mk-folder-item
  ><mk-file-item
    v-for="file in files"
    :key="file.path"
    :file="file"
    :depth="depth"
    @activated="onFile(file)"
  ></mk-file-item
></ul>
</template>
<script>

import mkFileItem from './FileItem.vue'
import mkFolderItem from './FolderItem.vue'
import { CatalogFolder } from './catalogItems.js'

export default {
  name: 'mkFolder',
  emits: ['fileSelected', 'folderSelected'],
  components: { mkFolderItem, mkFileItem },
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
    onFile(file) {
      this.$emit('fileSelected', file)
    },
    onFolder(folder) {
      this.$emit('folderSelected', folder);
    },
    items(isFolder) {
      const { folder } = this;
      return folder.contents? folder.contents.filter((el)=> {
        return (el instanceof CatalogFolder) === isFolder;
      }).sort((a,b)=> a.name.localeCompare(b.name)): [];
    },
  }
}
</script>

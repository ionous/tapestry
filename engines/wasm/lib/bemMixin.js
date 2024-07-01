export default function bemMixin(n) {
  if (!n) {
    throw new Error("missing name for bemMixin")
  }
  return {
    methods: {
      bemBlock(mod=false) {
        const blockName= n; // this.$options._componentTag -- doesnt exist in vue3
        const ar= [blockName];
        if (mod) {
          ar.push(blockName+"--"+mod);
        }
        return ar;
      },
      bemElem(el, mod=false) {
        const blockName= n; // ... could maybe use $options.file ... for now just require the name.
        const elName= blockName+ "__" + el;
        const ar= [elName];
        if (mod) {
          ar.push(elName+"--"+mod);
        }
        return ar;
      },
    }
  }
};

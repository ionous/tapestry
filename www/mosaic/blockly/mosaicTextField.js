'use strict';

import Blockly from 'blockly';

// custom text field to determine when text has a default value
// and use placeholder text instead.
export default class MosaicTextField extends Blockly.FieldTextInput {
  // called via class.fromJson()
  constructor(opt_value, opt_validator, opt_config) {
    // skip setup so that value is not set
    super(Blockly.Field.SKIP_SETUP, opt_validator, opt_config);
    this._defaultValue = opt_value;
    this.configure_(opt_config);
    this.setValue(''); // default value
    if (opt_validator) this.setValidator(opt_validator);
  }

  getDisplayText_() {
    return this.useDefaultDisplay() ? this._defaultValue : super.getDisplayText_();
  }

  useDefaultDisplay() {
    return !this.isBeingEdited_ && this.value_ === '';
  }

  // from field.js
  createTextElement_() {
    // https://developer.mozilla.org/en-US/docs/Web/SVG/Element/text
    super.createTextElement_();
    if (this.useDefaultDisplay()) {
      // index.css has .blocklyEditableText > text.blocklyText.placeholderText
      this.textElement_.classList.add('placeholderText');
    }
  }
}

// uses the default fromJson, classValidation, and default value ('')
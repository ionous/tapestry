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
  render_() {
    super.render_();
    // index.css has .blocklyEditableText > text.blocklyText.placeholderText
    // uses "fieldGroup_" instead of "textElement_" to try to match the hierarchy of the multiline text field.
    this.fieldGroup_.classList.toggle('placeholderText', this.useDefaultDisplay());
  }
}

// uses the default fromJson, classValidation, and default value ('')
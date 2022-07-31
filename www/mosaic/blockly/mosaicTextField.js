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
  }
  getDisplayText_() {
    return this.useDefaultDisplay() ? this._defaultValue : super.getDisplayText_();
  }

  useDefaultDisplay() {
    return !this.isBeingEdited_ && this.value_ === '';
  }

  // from field.js
  createTextElement_() {
    this.textElement_ = Blockly.utils.dom.createSvgElement(
        Blockly.utils.Svg.TEXT, {
          'class': !this.useDefaultDisplay() ? 'blocklyText' : 'blocklyText mosaicDefaultText',
        },
        this.fieldGroup_);
    if (this.getConstants().FIELD_TEXT_BASELINE_CENTER) {
      this.textElement_.setAttribute('dominant-baseline', 'central');
    }
    this.textContent_ = document.createTextNode('');
    this.textElement_.appendChild(this.textContent_);
  }
}

// uses the default fromJson, classValidation, and default value ('')
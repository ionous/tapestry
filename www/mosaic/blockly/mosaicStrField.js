'use strict';

import Blockly from 'blockly';

// for open strs with choices: allow authors to type their own custom text
// ( ex. determiners )
export default class MosaicStrField extends Blockly.FieldDropdown {
  // called via class.fromJson()
  // menuGenerator: an array of tuples: ex. [['are', '$ARE'], ...]
  // opt_config:{ name:'ARE_BEING', label:'are being',  options:the same menuGenerator array, type:'mosaic_str' }
  // it's not clear why FieldDropdown.fromJson() passes the options twice; we don't since we throw the value out anyway.
  constructor(_menuGenerator, opt_validator, opt_config) {
    const options = opt_config.options; // is it okay to not copy this? unclear.
    super(() => {
      // to avoid having options get 'trimmed' -- we have to pass a menuGenerator that's a function.
      // wish they had done this with a filter extension via fromJson or something via opt_config.
      return options;
    }, opt_validator, opt_config);
    this._constructed = true;
  }

  doValueUpdate_(newValue) {
    // FIX: some of our open str values (all of them probably) are being sent as their label not their key
    if (this._constructed && !newValue.startsWith('$')) {
      const options = this.menuGenerator_();
      for (let i=0; i < options.length; ++i) {
        const [ label, key ]  = options[i];
        if (label === newValue) {
          newValue = key;
          break
        }
      }
    }
    super.doValueUpdate_(newValue);
  }

  // note: called in many unexpected places to turn key into text and to manage an internal "selectedOption_" tuple.
  // ex. in the constructor used to find a default value for selectedOption, and then again in doValueUpdate_ due to setting that default value.
  // selectedOption mainly seems to exist for image dropdowns.
  // feels like it would have been better for value to have been an index, or '.options' be an object.
  // ( also a separate specialized class for image dropdowns would have been nice/r. )
  getOptions(_opt_useCache) {
    const val = this.getValue(); // initially: class.DEFAULT_VALUE ( null for FieldDropdown. )
    let options = this.menuGenerator_();
    if (val) {
      // did the author enter some custom text?
      if (!val.startsWith('$')) {
        // an entry with the label and key both set to the custom text.
        options = options.concat([[val, val]]);
      }
      // no clue why, blockly only allows strings as values.
      // ( it's validated in an internal function which we don't have access to )
      options = options.concat([['Set Custom Value...', '$']])
    }
    return options;
  }
  // menuItem is of type Blockly.MenuItem; 'menu' ( and this._menu ) appears to be null here.
  onItemSelected_(menu, menuItem) {
    const id = menuItem.getValue();
    if (id !== '$') {
      super.onItemSelected_(menu, menuItem);
    }
    else {
      this.showPromptEditor();
    }
  }

  // pulled from field_textinput.js
  // TODO: some sort of inline text editor?
  showPromptEditor() {
    // getText returns the label of our currently selected value
    Blockly.dialog.prompt(Blockly.Msg['CHANGE_VALUE_TITLE'], this.getText(), (text) => {
      // Text is null if user pressed cancel button.
      if (text !== null) {
        // setValue for the dropdown uses the "key" value ex. "$BEEP"
        // if its custom text, then the key is the requested text value
        this.setValue(String(text));
      }
    });
  }
}

MosaicStrField.fromJson = function(options) {
  // replace "%{BKY_MY_MSG}" with the value in Msg['MY_MSG'].
  // var value = Blockly.utils.replaceMessageReferences(options['value']);
  return new MosaicStrField(/*options['options']*/null, undefined, options);
};


// no validation. im not even sure why blockly has validation for dropdowns in the first place.
// the options are already validated so the only way their dropdown value could fail is their own code going wrong.
MosaicStrField.prototype.doClassValidation_ = function(opt_newValue) {};
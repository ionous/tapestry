"use strict";

import Blockly from "blockly";

// a clone of MosaicTextField for placeholder text
// ( can this be a mixin instead? can a mixin intercept constructors? )
// also fixes svg sizing issues on load.
export default class MosaicMultilineField extends Blockly.FieldMultilineInput {
  // called via class.fromJson()
  constructor(opt_value, opt_validator, opt_config) {
    // skip setup so that value is not set
    super(Blockly.Field.SKIP_SETUP, opt_validator, opt_config);
    this._defaultValue = opt_value; // in tapestry: this is the name of the field; the real value gets loaded later.
    this.configure_(opt_config);
    this.setValue(""); // default value
    if (opt_validator) this.setValidator(opt_validator);
  }

  getDisplayText_() {
    return this.useDefaultDisplay()
      ? this._defaultValue
      : super.getDisplayText_();
  }

  useDefaultDisplay() {
    return !this.isBeingEdited_ && this.value_ === "";
  }

  render_() {
    super.render_();
    // index.css has .blocklyEditableText > text.blocklyText.placeholderText
    this.textGroup_.classList.toggle('placeholderText', this.useDefaultDisplay());
  }
  /**
   * Updates the size of the field based on the text.
   * @protected
   */
  updateSize_() {
    const dom = Blockly.utils.dom;
    const fontSize = this.getConstants().FIELD_TEXT_FONTSIZE;
    const fontWeight = this.getConstants().FIELD_TEXT_FONTWEIGHT;
    const fontFamily = this.getConstants().FIELD_TEXT_FONTFAMILY;
    const scale = this.sourceBlock_.workspace.getScale();
    //
    const nodes = this.textGroup_.childNodes;
    let totalWidth = 0;
    let totalHeight = 0;
    for (let i = 0; i < nodes.length; i++) {
      const tspan = /** @type {!Element} */ (nodes[i]);
      // mod-stravis: https://github.com/google/blockly/issues/6071
      const textWidth = dom.getTextWidth(tspan) || (dom.getFastTextWidth(tspan, fontSize, fontWeight, fontFamily));
      if (textWidth > totalWidth) {
        totalWidth = textWidth;
      }
      totalHeight +=
        this.getConstants().FIELD_TEXT_HEIGHT +
        (i > 0 ? this.getConstants().FIELD_BORDER_RECT_Y_PADDING : 0);
    }
    if (this.isBeingEdited_) {
      // The default width is based on the longest line in the display text,
      // but when it's being edited, width should be calculated based on the
      // absolute longest line, even if it would be truncated after editing.
      // Otherwise we would get wrong editor width when there are more
      // lines than this.maxLines_.
      const actualEditorLines = this.value_.split("\n");
      const dummyTextElement = dom.createSvgElement(Blockly.utils.Svg.TEXT, {
        class: "blocklyText blocklyMultilineText",
      });

      for (let i = 0; i < actualEditorLines.length; i++) {
        if (actualEditorLines[i].length > this.maxDisplayLength) {
          actualEditorLines[i] = actualEditorLines[i].substring(
            0,
            this.maxDisplayLength
          );
        }
        dummyTextElement.textContent = actualEditorLines[i];
        const lineWidth = dom.getFastTextWidth(
          dummyTextElement,
          fontSize,
          fontWeight,
          fontFamily
        );
        if (lineWidth > totalWidth) {
          totalWidth = lineWidth;
        }
      }

      const scrollbarWidth =
        this.htmlInput_.offsetWidth - this.htmlInput_.clientWidth;
      totalWidth += scrollbarWidth;
    }
    if (this.borderRect_) {
      // inflate the border rect
      totalWidth += this.getConstants().FIELD_BORDER_RECT_X_PADDING * 2;
      totalHeight +=  this.getConstants().FIELD_BORDER_RECT_Y_PADDING * 2;
      // mod-stravis: --> these get overriden by positionBorderRect_
      //this.borderRect_.setAttribute("width", totalWidth);
      //this.borderRect_.setAttribute("height", totalHeight);
    }
    this.size_.width = totalWidth;
    this.size_.height = totalHeight;

    this.positionBorderRect_();
  }
}

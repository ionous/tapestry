const jsonData = [
  {
    "type": "bool_value",
    "output": [
      "bool_value",
      "bool_eval",
      "literal_value"
    ],
    "message0": "bool",
    "colour": "%{BKY_LOGIC_HUE}",
    "tooltip": "Specify an explicit true or false.",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "bool"
          },
          {
            "type": "field_input"
          },
          {
            "name": "BOOL",
            "type": "input_dummy"
          }
        ],
        [
          {
            "type": "field_label",
            "text": "class"
          },
          {
            "type": "field_input"
          },
          {
            "name": "CLASS",
            "type": "input_dummy",
            "optional": true
          }
        ]
      ]
    }
  },
  {
    "type": "_bool_value_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "bool"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "class"
      },
      {
        "type": "field_checkbox",
        "name": "class_edit"
      },
      {
        "name": "CLASS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5"
  },
  {
    "type": "field_value",
    "output": [
      "field_value"
    ],
    "message0": "field",
    "colour": "%{BKY_COLOUR_HUE}",
    "tooltip": "A fixed value of a record.",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "field"
          },
          {
            "type": "field_input"
          },
          {
            "name": "FIELD",
            "type": "input_dummy"
          }
        ],
        [
          {
            "type": "field_label",
            "text": "value"
          },
          {
            "name": "VALUE",
            "type": "input_value",
            "check": [
              "literal_value"
            ]
          }
        ]
      ]
    }
  },
  {
    "type": "_field_value_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "field"
      },
      {
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2"
  },
  {
    "type": "field_values",
    "output": [
      "field_values",
      "literal_value"
    ],
    "message0": "fields",
    "colour": "%{BKY_MATH_HUE}",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "contains"
          },
          {
            "name": "CONTAINS",
            "type": "input_value",
            "check": [
              "field_value"
            ],
            "shadow": "field_value",
            "repeats": true
          }
        ]
      ]
    }
  },
  {
    "type": "_field_values_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "fields"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "contains"
      },
      {
        "type": "field_number",
        "min": 1,
        "precision": 1,
        "name": "contains_edit"
      },
      {
        "name": "CONTAINS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5"
  },
  {
    "type": "num_value",
    "output": [
      "num_value",
      "number_eval",
      "literal_value"
    ],
    "message0": "num",
    "colour": "%{BKY_MATH_HUE}",
    "tooltip": "Specify a particular number.",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "num"
          },
          {
            "type": "field_number"
          },
          {
            "name": "NUM",
            "type": "input_dummy"
          }
        ],
        [
          {
            "type": "field_label",
            "text": "class"
          },
          {
            "type": "field_input"
          },
          {
            "name": "CLASS",
            "type": "input_dummy",
            "optional": true
          }
        ]
      ]
    }
  },
  {
    "type": "_num_value_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "num"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "class"
      },
      {
        "type": "field_checkbox",
        "name": "class_edit"
      },
      {
        "name": "CLASS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5"
  },
  {
    "type": "num_values",
    "output": [
      "num_values",
      "num_list_eval",
      "literal_value"
    ],
    "message0": "nums",
    "colour": "%{BKY_MATH_HUE}",
    "tooltip": "Number List: Specify a list of numbers.",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "values"
          },
          {
            "type": "field_number"
          },
          {
            "name": "VALUES",
            "type": "input_dummy",
            "repeats": true
          }
        ],
        [
          {
            "type": "field_label",
            "text": "class"
          },
          {
            "type": "field_input"
          },
          {
            "name": "CLASS",
            "type": "input_dummy",
            "optional": true
          }
        ]
      ]
    }
  },
  {
    "type": "_num_values_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "nums"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "values"
      },
      {
        "type": "field_number",
        "min": 1,
        "precision": 1,
        "name": "values_edit"
      },
      {
        "name": "VALUES",
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "class"
      },
      {
        "type": "field_checkbox",
        "name": "class_edit"
      },
      {
        "name": "CLASS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5%6%7%8"
  },
  {
    "type": "record_value",
    "output": [
      "record_value",
      "record_eval",
      "literal_value"
    ],
    "message0": "rec",
    "colour": "%{BKY_LISTS_HUE}",
    "tooltip": "Specify a record composed of literal values.",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "kind"
          },
          {
            "type": "field_input"
          },
          {
            "name": "KIND",
            "type": "input_dummy"
          }
        ],
        [
          {
            "type": "field_label",
            "text": "fields"
          },
          {
            "name": "FIELDS",
            "type": "input_value",
            "check": [
              "field_value"
            ],
            "shadow": "field_value",
            "repeats": true
          }
        ]
      ]
    }
  },
  {
    "type": "_record_value_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "rec"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "fields"
      },
      {
        "type": "field_number",
        "min": 1,
        "precision": 1,
        "name": "fields_edit"
      },
      {
        "name": "FIELDS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5"
  },
  {
    "type": "record_values",
    "output": [
      "record_values",
      "record_list_eval",
      "literal_value"
    ],
    "message0": "recs",
    "colour": "%{BKY_LISTS_HUE}",
    "tooltip": "Specify a series of records, all of the same kind.",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "kind"
          },
          {
            "type": "field_input"
          },
          {
            "name": "KIND",
            "type": "input_dummy"
          }
        ],
        [
          {
            "type": "field_label",
            "text": "containing"
          },
          {
            "name": "ELS",
            "type": "input_value",
            "check": [
              "field_values"
            ],
            "shadow": "field_values",
            "repeats": true
          }
        ]
      ]
    }
  },
  {
    "type": "_record_values_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "recs"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "containing"
      },
      {
        "type": "field_number",
        "min": 1,
        "precision": 1,
        "name": "els_edit"
      },
      {
        "name": "ELS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5"
  },
  {
    "type": "test_flow",
    "output": [
      "test_flow",
      "test_slot"
    ],
    "message0": "flow",
    "colour": "%{BKY_COLOUR_HUE}",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "slot"
          },
          {
            "name": "SLOT",
            "type": "input_value",
            "check": [
              "test_slot"
            ],
            "optional": true
          }
        ],
        [
          {
            "type": "field_label",
            "text": "txt"
          },
          {
            "type": "field_input"
          },
          {
            "name": "TXT",
            "type": "input_dummy",
            "optional": true
          }
        ],
        [
          {
            "type": "field_label",
            "text": "num"
          },
          {
            "type": "field_number"
          },
          {
            "name": "NUM",
            "type": "input_dummy",
            "optional": true
          }
        ],
        [
          {
            "type": "field_label",
            "text": "bool"
          },
          {
            "type": "field_input"
          },
          {
            "name": "BOOL",
            "type": "input_dummy",
            "optional": true
          }
        ],
        [
          {
            "type": "field_label",
            "text": "swap"
          },
          {
            "type": "field_dropdown",
            "options": [
              [
                "a",
                "test_flow"
              ],
              [
                "b",
                "test_slot"
              ],
              [
                "c",
                "test_txt"
              ]
            ]
          },
          {
            "name": "SWAP",
            "type": "input_value",
            "check": [
              "test_flow",
              "test_slot",
              "test_txt"
            ],
            "optional": true
          }
        ],
        [
          {
            "type": "field_label",
            "text": "slots"
          },
          {
            "name": "SLOTS",
            "type": "input_value",
            "check": [
              "test_slot"
            ],
            "optional": true,
            "repeats": true
          }
        ]
      ]
    }
  },
  {
    "type": "_test_flow_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "flow"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "slot"
      },
      {
        "type": "field_checkbox",
        "name": "slot_edit"
      },
      {
        "name": "SLOT",
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "txt"
      },
      {
        "type": "field_checkbox",
        "name": "txt_edit"
      },
      {
        "name": "TXT",
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "num"
      },
      {
        "type": "field_checkbox",
        "name": "num_edit"
      },
      {
        "name": "NUM",
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "bool"
      },
      {
        "type": "field_checkbox",
        "name": "bool_edit"
      },
      {
        "name": "BOOL",
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "swap"
      },
      {
        "type": "field_checkbox",
        "name": "swap_edit"
      },
      {
        "name": "SWAP",
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "slots"
      },
      {
        "type": "field_number",
        "min": 0,
        "precision": 1,
        "name": "slots_edit"
      },
      {
        "name": "SLOTS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5%6%7%8%9%10%11%12%13%14%15%16%17%18%19%20"
  },
  {
    "type": "test_str",
    "output": [
      "test_str",
      "test_slot"
    ],
    "message0": "test_str",
    "colour": "%{BKY_COLOUR_HUE}",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_input"
          },
          {
            "name": "TEST_STR",
            "type": "input_dummy"
          }
        ]
      ]
    }
  },
  {
    "type": "text_value",
    "output": [
      "text_value",
      "text_eval",
      "literal_value"
    ],
    "message0": "txt",
    "colour": "%{BKY_TEXTS_HUE}",
    "tooltip": "Specify a small bit of text.",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "text"
          },
          {
            "type": "field_input"
          },
          {
            "name": "TEXT",
            "type": "input_dummy"
          }
        ],
        [
          {
            "type": "field_label",
            "text": "class"
          },
          {
            "type": "field_input"
          },
          {
            "name": "CLASS",
            "type": "input_dummy",
            "optional": true
          }
        ]
      ]
    }
  },
  {
    "type": "_text_value_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "txt"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "class"
      },
      {
        "type": "field_checkbox",
        "name": "class_edit"
      },
      {
        "name": "CLASS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5"
  },
  {
    "type": "text_values",
    "output": [
      "text_values",
      "text_list_eval",
      "literal_value"
    ],
    "message0": "txts",
    "colour": "%{BKY_TEXTS_HUE}",
    "tooltip": "Text List: Specifies a set of text values.",
    "extensions": [
      "tapestry_generic_mixin",
      "tapestry_generic_extension"
    ],
    "mutator": "tapestry_generic_mutation",
    "customData": {
      "blockDef": [
        [
          {
            "type": "field_label",
            "text": "values"
          },
          {
            "type": "field_input"
          },
          {
            "name": "VALUES",
            "type": "input_dummy",
            "repeats": true
          }
        ],
        [
          {
            "type": "field_label",
            "text": "class"
          },
          {
            "type": "field_input"
          },
          {
            "name": "CLASS",
            "type": "input_dummy",
            "optional": true
          }
        ]
      ]
    }
  },
  {
    "type": "_text_values_mutator",
    "style": "logic_blocks",
    "inputsInline": false,
    "args0": [
      {
        "type": "field_label",
        "text": "txts"
      },
      {
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "values"
      },
      {
        "type": "field_number",
        "min": 1,
        "precision": 1,
        "name": "values_edit"
      },
      {
        "name": "VALUES",
        "type": "input_dummy"
      },
      {
        "type": "field_label",
        "text": "class"
      },
      {
        "type": "field_checkbox",
        "name": "class_edit"
      },
      {
        "name": "CLASS",
        "type": "input_dummy"
      }
    ],
    "message0": "%1%2%3%4%5%6%7%8"
  }
]

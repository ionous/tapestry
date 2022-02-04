const jsonData = [
          {
            "type": "abstract_action",
            "output": [
              "abstract_action"
            ],
            "message0": "abstract_action",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_abstract_action_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "nothing",
                        "$NOTHING"
                      ]
                    ]
                  },
                  {
                    "name": "ABSTRACT_ACTION",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "action",
            "output": [
              "action",
              "scanner_maker"
            ],
            "message0": "as",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "makes a parser scanner producing a script defined action.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_action_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "action"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ACTION",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_action_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "as"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "action_context",
            "output": [
              "action_context"
            ],
            "message0": "action_context",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_action_context_mutator",
              "shapeDef": [
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
                ]
              ]
            }
          },
          {
            "type": "_action_context_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "action_context"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_action_decl_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "action_decl",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Declare an activity: Activities help actors perform tasks: for instance, picking up or dropping items.  Activities involve either the player or an npc and possibly one or two other objects.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_action_decl_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "act"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "EVENT",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "acting"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ACTION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "action_params"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one or more objects",
                        "$COMMON"
                      ],
                      [
                        "two similar objects",
                        "$DUAL"
                      ],
                      [
                        "nothing",
                        "$NONE"
                      ]
                    ],
                    "swaps": {
                      "$COMMON": "common_action",
                      "$DUAL": "paired_action",
                      "$NONE": "abstract_action"
                    }
                  },
                  {
                    "name": "ACTION_PARAMS",
                    "type": "input_value",
                    "checks": [
                      "common_action",
                      "paired_action",
                      "abstract_action"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "action_decl",
            "output": [
              "action_decl"
            ],
            "message0": "action_decl",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Declare an activity: Activities help actors perform tasks: for instance, picking up or dropping items.  Activities involve either the player or an npc and possibly one or two other objects.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_action_decl_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "act"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "EVENT",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "acting"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ACTION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "action_params"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one or more objects",
                        "$COMMON"
                      ],
                      [
                        "two similar objects",
                        "$DUAL"
                      ],
                      [
                        "nothing",
                        "$NONE"
                      ]
                    ],
                    "swaps": {
                      "$COMMON": "common_action",
                      "$DUAL": "paired_action",
                      "$NONE": "abstract_action"
                    }
                  },
                  {
                    "name": "ACTION_PARAMS",
                    "type": "input_value",
                    "checks": [
                      "common_action",
                      "paired_action",
                      "abstract_action"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_action_decl_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "action_decl"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "action_name",
            "output": [
              "action_name"
            ],
            "message0": "action_name",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_action_name_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ACTION_NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "action_params",
            "output": [
              "action_params"
            ],
            "message0": "action_params",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_action_params_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one or more objects",
                        "$COMMON"
                      ],
                      [
                        "two similar objects",
                        "$DUAL"
                      ],
                      [
                        "nothing",
                        "$NONE"
                      ]
                    ],
                    "swaps": {
                      "$COMMON": "common_action",
                      "$DUAL": "paired_action",
                      "$NONE": "abstract_action"
                    }
                  },
                  {
                    "name": "ACTION_PARAMS",
                    "type": "input_value",
                    "checks": [
                      "common_action",
                      "paired_action",
                      "abstract_action"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_activity_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "act",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_activity_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "exe"
                  },
                  {
                    "name": "EXE",
                    "type": "input_statement",
                    "checks": [
                      "_execute_stack"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "activity",
            "output": [
              "activity"
            ],
            "message0": "act",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_activity_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "exe"
                  },
                  {
                    "name": "EXE",
                    "type": "input_statement",
                    "checks": [
                      "_execute_stack"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_activity_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "act"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "affinity",
            "output": [
              "affinity"
            ],
            "message0": "affinity",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Indicates storage for fields and other properties.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_affinity_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "bool",
                        "$BOOL"
                      ],
                      [
                        "number",
                        "$NUMBER"
                      ],
                      [
                        "num_list",
                        "$NUM_LIST"
                      ],
                      [
                        "text",
                        "$TEXT"
                      ],
                      [
                        "text_list",
                        "$TEXT_LIST"
                      ],
                      [
                        "record",
                        "$RECORD"
                      ],
                      [
                        "record_list",
                        "$RECORD_LIST"
                      ]
                    ]
                  },
                  {
                    "name": "AFFINITY",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "alias",
            "output": [
              "alias",
              "grammar_maker"
            ],
            "message0": "alias",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "allows the user to refer to a noun by one or more other terms.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_alias_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "names"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAMES",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as_noun"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "AS_NOUN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_alias_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "alias"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "names"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "names_edit"
              },
              {
                "name": "NAMES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "all_of",
            "output": [
              "all_of",
              "scanner_maker"
            ],
            "message0": "all_of",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "makes a parser scanner.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_all_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "series"
                  },
                  {
                    "name": "SERIES",
                    "type": "input_value",
                    "checks": [
                      "scanner_maker"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_all_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "all_of"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "series"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "series_edit"
              },
              {
                "name": "SERIES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "all_true",
            "output": [
              "all_true",
              "bool_eval"
            ],
            "message0": "all_true",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Returns true if all of the evaluations are true.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_all_true_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test"
                  },
                  {
                    "name": "TEST",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_all_true_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "all_true"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "test"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "test_edit"
              },
              {
                "name": "TEST",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "always",
            "output": [
              "always",
              "bool_eval"
            ],
            "message0": "always",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Returns true.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_always_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_always_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "always"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "ana",
            "output": [
              "ana"
            ],
            "message0": "ana",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_ana_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "a",
                        "$A"
                      ],
                      [
                        "an",
                        "$AN"
                      ]
                    ]
                  },
                  {
                    "name": "ANA",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "any_of",
            "output": [
              "any_of",
              "scanner_maker"
            ],
            "message0": "any_of",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "makes a parser scanner.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_any_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "options"
                  },
                  {
                    "name": "OPTIONS",
                    "type": "input_value",
                    "checks": [
                      "scanner_maker"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_any_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "any_of"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "options"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "options_edit"
              },
              {
                "name": "OPTIONS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "any_true",
            "output": [
              "any_true",
              "bool_eval"
            ],
            "message0": "any_true",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Returns true if any of the evaluations are true.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_any_true_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test"
                  },
                  {
                    "name": "TEST",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_any_true_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "any_true"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "test"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "test_edit"
              },
              {
                "name": "TEST",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "are_an",
            "output": [
              "are_an"
            ],
            "message0": "are_an",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_are_an_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "are a",
                        "$AREA"
                      ],
                      [
                        "are an",
                        "$AREAN"
                      ],
                      [
                        "is",
                        "$IS"
                      ],
                      [
                        "is a",
                        "$ISA"
                      ],
                      [
                        "is an",
                        "$ISAN"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_AN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "are_being",
            "output": [
              "are_being"
            ],
            "message0": "are_being",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_are_being_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "is",
                        "$IS"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_BEING",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "are_either",
            "output": [
              "are_either"
            ],
            "message0": "are_either",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_are_either_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "can be",
                        "$CANBE"
                      ],
                      [
                        "are either",
                        "$EITHER"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_EITHER",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "arg",
            "output": [
              "arg"
            ],
            "message0": "arg",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Runtime version of argument.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_arg_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_arg_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "arg"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "as_num",
            "output": [
              "as_num",
              "list_iterator"
            ],
            "message0": "as_num",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Define the name of a number variable.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_as_num_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_as_num_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "as_num"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "as_rec",
            "output": [
              "as_rec",
              "list_iterator"
            ],
            "message0": "as_rec",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Define the name of a record variable.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_as_rec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_as_rec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "as_rec"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "as_txt",
            "output": [
              "as_txt",
              "list_iterator"
            ],
            "message0": "as_txt",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Define the name of a text variable.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_as_txt_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_as_txt_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "as_txt"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "aspect",
            "output": [
              "aspect"
            ],
            "message0": "aspect",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_aspect_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASPECT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "aspect_property",
            "output": [
              "aspect_property",
              "property_slot"
            ],
            "message0": "aspect_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_aspect_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "aspect"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASPECT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_aspect_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "aspect_property"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_aspect_traits_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "aspect_traits",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add traits to an aspect",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_aspect_traits_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "aspect"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASPECT",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait_phrase"
                  },
                  {
                    "name": "TRAIT_PHRASE",
                    "type": "input_value",
                    "checks": [
                      "trait_phrase"
                    ],
                    "shadow": "trait_phrase"
                  }
                ]
              ]
            }
          },
          {
            "type": "aspect_traits",
            "output": [
              "aspect_traits"
            ],
            "message0": "aspect_traits",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add traits to an aspect",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_aspect_traits_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "aspect"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASPECT",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait_phrase"
                  },
                  {
                    "name": "TRAIT_PHRASE",
                    "type": "input_value",
                    "checks": [
                      "trait_phrase"
                    ],
                    "shadow": "trait_phrase"
                  }
                ]
              ]
            }
          },
          {
            "type": "_aspect_traits_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "aspect_traits"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_assign_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "let",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Assigns a variable to a value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_assign_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "be"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "assign",
            "output": [
              "assign"
            ],
            "message0": "let",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Assigns a variable to a value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_assign_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "be"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_assign_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "let"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "at_least",
            "output": [
              "at_least",
              "comparator"
            ],
            "message0": "at_least",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The first value is greater than or equal to the second value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_at_least_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_at_least_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "at_least"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "at_most",
            "output": [
              "at_most",
              "comparator"
            ],
            "message0": "at_most",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The first value is less than or equal to the second value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_at_most_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_at_most_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "at_most"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_blankline_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "p",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add a single blank line following some text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_blankline_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "blankline",
            "output": [
              "blankline"
            ],
            "message0": "p",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add a single blank line following some text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_blankline_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_blankline_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "p"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "bool",
            "output": [
              "bool"
            ],
            "message0": "bool",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_bool_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
                  },
                  {
                    "name": "BOOL",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "bool_property",
            "output": [
              "bool_property",
              "property_slot"
            ],
            "message0": "bool_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_bool_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "initially"
                  },
                  {
                    "name": "INITIALLY",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_bool_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "bool_property"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "initially"
              },
              {
                "type": "field_checkbox",
                "name": "initially_edit"
              },
              {
                "name": "INITIALLY",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
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
              "mui": "_bool_value_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "bool"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
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
            "type": "bracket_text",
            "output": [
              "bracket_text",
              "text_eval"
            ],
            "message0": "brackets",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Sandwiches text printed during a block and puts them inside parenthesis '()'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_bracket_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_bracket_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "brackets"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_break_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "break",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "In a repeating loop, exit the loop.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_break_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "break",
            "output": [
              "break"
            ],
            "message0": "break",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "In a repeating loop, exit the loop.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_break_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_break_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "break"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "buffer_text",
            "output": [
              "buffer_text",
              "text_eval"
            ],
            "message0": "buffers",
            "colour": "%{BKY_TEXTS_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_buffer_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_buffer_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "buffers"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "call_args",
            "output": [
              "call_args"
            ],
            "message0": "args",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Pattern arguments.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_args_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "args"
                  },
                  {
                    "name": "ARGS",
                    "type": "input_value",
                    "checks": [
                      "arg"
                    ],
                    "shadow": "arg",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_call_args_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "args"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "args"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "args_edit"
              },
              {
                "name": "ARGS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "call_cycle",
            "output": [
              "call_cycle",
              "text_eval"
            ],
            "message0": "cycle",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Runtime version of cycle_text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_cycle_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "over"
                  },
                  {
                    "name": "PARTS",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_call_cycle_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "cycle"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "over"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "parts_edit"
              },
              {
                "name": "PARTS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_call_pattern_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "determine",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Executes a pattern, and potentially returns a value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_pattern_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "pattern"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PATTERN",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "args"
                  },
                  {
                    "name": "ARGUMENTS",
                    "type": "input_value",
                    "checks": [
                      "call_args"
                    ],
                    "shadow": "call_args"
                  }
                ]
              ]
            }
          },
          {
            "type": "call_pattern",
            "output": [
              "call_pattern",
              "bool_eval",
              "number_eval",
              "text_eval",
              "record_eval",
              "num_list_eval",
              "text_list_eval",
              "record_list_eval"
            ],
            "message0": "determine",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Executes a pattern, and potentially returns a value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_pattern_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "pattern"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PATTERN",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "args"
                  },
                  {
                    "name": "ARGUMENTS",
                    "type": "input_value",
                    "checks": [
                      "call_args"
                    ],
                    "shadow": "call_args"
                  }
                ]
              ]
            }
          },
          {
            "type": "_call_pattern_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "determine"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_call_send_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "send",
            "colour": "%{BKY_LOGIC_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_send_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "path"
                  },
                  {
                    "name": "PATH",
                    "type": "input_value",
                    "checks": [
                      "text_list_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "event"
                  },
                  {
                    "name": "EVENT",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "call_send",
            "output": [
              "call_send",
              "bool_eval"
            ],
            "message0": "send",
            "colour": "%{BKY_LOGIC_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_send_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "path"
                  },
                  {
                    "name": "PATH",
                    "type": "input_value",
                    "checks": [
                      "text_list_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "event"
                  },
                  {
                    "name": "EVENT",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_call_send_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "send"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "call_shuffle",
            "output": [
              "call_shuffle",
              "text_eval"
            ],
            "message0": "shuffle",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Runtime version of shuffle_text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_shuffle_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "over"
                  },
                  {
                    "name": "PARTS",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_call_shuffle_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "shuffle"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "over"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "parts_edit"
              },
              {
                "name": "PARTS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "call_terminal",
            "output": [
              "call_terminal",
              "text_eval"
            ],
            "message0": "stopping",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Runtime version of stopping_text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_terminal_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "over"
                  },
                  {
                    "name": "PARTS",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_call_terminal_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "stopping"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "over"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "parts_edit"
              },
              {
                "name": "PARTS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "call_trigger",
            "output": [
              "call_trigger",
              "bool_eval"
            ],
            "message0": "trigger",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Runtime version of count_of.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_call_trigger_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "on"
                  },
                  {
                    "name": "TRIGGER",
                    "type": "input_value",
                    "checks": [
                      "trigger"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "num"
                  },
                  {
                    "name": "NUM",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_call_trigger_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "trigger"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "capitalize",
            "output": [
              "capitalize",
              "text_eval"
            ],
            "message0": "capitalize",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns new text, with the first letter turned into uppercase.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_capitalize_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_capitalize_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "capitalize"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_certainties_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "certainties",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Give a kind a trait",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_certainties_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "plural_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL_KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "are_being"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "is",
                        "$IS"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_BEING",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "certainty"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "usually",
                        "$USUALLY"
                      ],
                      [
                        "always",
                        "$ALWAYS"
                      ],
                      [
                        "seldom",
                        "$SELDOM"
                      ],
                      [
                        "never",
                        "$NEVER"
                      ]
                    ]
                  },
                  {
                    "name": "CERTAINTY",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TRAIT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "certainties",
            "output": [
              "certainties"
            ],
            "message0": "certainties",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Give a kind a trait",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_certainties_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "plural_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL_KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "are_being"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "is",
                        "$IS"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_BEING",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "certainty"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "usually",
                        "$USUALLY"
                      ],
                      [
                        "always",
                        "$ALWAYS"
                      ],
                      [
                        "seldom",
                        "$SELDOM"
                      ],
                      [
                        "never",
                        "$NEVER"
                      ]
                    ]
                  },
                  {
                    "name": "CERTAINTY",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TRAIT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_certainties_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "certainties"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "certainty",
            "output": [
              "certainty"
            ],
            "message0": "certainty",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Whether an trait applies to a kind of noun.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_certainty_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "usually",
                        "$USUALLY"
                      ],
                      [
                        "always",
                        "$ALWAYS"
                      ],
                      [
                        "seldom",
                        "$SELDOM"
                      ],
                      [
                        "never",
                        "$NEVER"
                      ]
                    ]
                  },
                  {
                    "name": "CERTAINTY",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "choice_spec",
            "output": [
              "choice_spec"
            ],
            "message0": "pick",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choice_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "label"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LABEL",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "type"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_choice_spec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "pick"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "label"
              },
              {
                "type": "field_checkbox",
                "name": "label_edit"
              },
              {
                "name": "LABEL",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "type"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "_choose_action_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "if",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "An if statement.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_action_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "if"
                  },
                  {
                    "name": "IF",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "choose_action",
            "output": [
              "choose_action",
              "brancher"
            ],
            "message0": "if",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "An if statement.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_action_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "if"
                  },
                  {
                    "name": "IF",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_choose_action_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "if"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "else"
              },
              {
                "type": "field_checkbox",
                "name": "else_edit"
              },
              {
                "name": "ELSE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "choose_more",
            "output": [
              "choose_more",
              "brancher"
            ],
            "message0": "else_if",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_more_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "if"
                  },
                  {
                    "name": "IF",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_choose_more_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "else_if"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "else"
              },
              {
                "type": "field_checkbox",
                "name": "else_edit"
              },
              {
                "name": "ELSE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "choose_more_value",
            "output": [
              "choose_more_value",
              "brancher"
            ],
            "message0": "else_if",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_more_value_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "assign"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASSIGN",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "FILTER",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_choose_more_value_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "else_if"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "else"
              },
              {
                "type": "field_checkbox",
                "name": "else_edit"
              },
              {
                "name": "ELSE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "choose_nothing_else",
            "output": [
              "choose_nothing_else",
              "brancher"
            ],
            "message0": "else_do",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_nothing_else_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_choose_nothing_else_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "else_do"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "choose_num",
            "output": [
              "choose_num",
              "number_eval"
            ],
            "message0": "num",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Pick one of two numbers based on a boolean test.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_num_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "if"
                  },
                  {
                    "name": "IF",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "then"
                  },
                  {
                    "name": "TRUE",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "FALSE",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_choose_num_mutator",
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
                "text": "else"
              },
              {
                "type": "field_checkbox",
                "name": "false_edit"
              },
              {
                "name": "FALSE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "choose_text",
            "output": [
              "choose_text",
              "text_eval"
            ],
            "message0": "txt",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Pick one of two strings based on a boolean test.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "if"
                  },
                  {
                    "name": "IF",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "then"
                  },
                  {
                    "name": "TRUE",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "FALSE",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_choose_text_mutator",
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
                "text": "else"
              },
              {
                "type": "field_checkbox",
                "name": "false_edit"
              },
              {
                "name": "FALSE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_choose_value_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "if",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "An if statement with local assignment.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_value_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "assign"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASSIGN",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "FILTER",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "choose_value",
            "output": [
              "choose_value",
              "brancher"
            ],
            "message0": "if",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "An if statement with local assignment.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_choose_value_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "assign"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASSIGN",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "FILTER",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_choose_value_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "if"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "else"
              },
              {
                "type": "field_checkbox",
                "name": "else_edit"
              },
              {
                "name": "ELSE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "comma_text",
            "output": [
              "comma_text",
              "text_eval"
            ],
            "message0": "commas",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Separates words with commas, and 'and'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_comma_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_comma_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "commas"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_comment_stack",
            "nextStatement": [
              "_story_statement_stack",
              "_execute_stack"
            ],
            "previousStatement": [
              "_story_statement_stack",
              "_execute_stack"
            ],
            "message0": "comment",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add a note.\nInformation about the story for you and other authors.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_comment_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "comment"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LINES",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "comment",
            "output": [
              "comment"
            ],
            "message0": "comment",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add a note.\nInformation about the story for you and other authors.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_comment_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "comment"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LINES",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_comment_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "comment"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "common_action",
            "output": [
              "common_action"
            ],
            "message0": "common_action",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_common_action_mutator",
              "shapeDef": [
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
                    "text": "action_context"
                  },
                  {
                    "name": "ACTION_CONTEXT",
                    "type": "input_value",
                    "checks": [
                      "action_context"
                    ],
                    "shadow": "action_context",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_common_action_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "common_action"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "action_context"
              },
              {
                "type": "field_checkbox",
                "name": "action_context_edit"
              },
              {
                "name": "ACTION_CONTEXT",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "compare_num",
            "output": [
              "compare_num",
              "bool_eval"
            ],
            "message0": "cmp",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "True if eq,ne,gt,lt,ge,le two numbers.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_compare_num_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "a"
                  },
                  {
                    "name": "A",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "is"
                  },
                  {
                    "name": "IS",
                    "type": "input_value",
                    "checks": [
                      "comparator"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "num"
                  },
                  {
                    "name": "B",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_compare_num_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "cmp"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "compare_text",
            "output": [
              "compare_text",
              "bool_eval"
            ],
            "message0": "cmp",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "True if eq,ne,gt,lt,ge,le two strings ( lexical. ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_compare_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "a"
                  },
                  {
                    "name": "A",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "is"
                  },
                  {
                    "name": "IS",
                    "type": "input_value",
                    "checks": [
                      "comparator"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "txt"
                  },
                  {
                    "name": "B",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_compare_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "cmp"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "count_of",
            "output": [
              "count_of",
              "bool_eval"
            ],
            "message0": "count_of",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "A guard which returns true based on a counter.\nCounters start at zero and are incremented every time the guard gets checked.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_count_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "trigger"
                  },
                  {
                    "name": "TRIGGER",
                    "type": "input_value",
                    "checks": [
                      "trigger"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "num"
                  },
                  {
                    "name": "NUM",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_count_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "count_of"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "cycle_text",
            "output": [
              "cycle_text",
              "text_eval"
            ],
            "message0": "cycle_text",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "When called multiple times, returns each of its inputs in turn.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_cycle_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "parts"
                  },
                  {
                    "name": "PARTS",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_cycle_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "cycle_text"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "parts"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "parts_edit"
              },
              {
                "name": "PARTS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_debug_log_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "log",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Debug log.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_debug_log_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "value"
                  },
                  {
                    "name": "VALUE",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "note",
                        "$NOTE"
                      ],
                      [
                        "to_do",
                        "$TO_DO"
                      ],
                      [
                        "fix",
                        "$FIX"
                      ],
                      [
                        "info",
                        "$INFO"
                      ],
                      [
                        "warning",
                        "$WARNING"
                      ],
                      [
                        "error",
                        "$ERROR"
                      ]
                    ]
                  },
                  {
                    "name": "LOG_LEVEL",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "debug_log",
            "output": [
              "debug_log"
            ],
            "message0": "log",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Debug log.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_debug_log_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "value"
                  },
                  {
                    "name": "VALUE",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "note",
                        "$NOTE"
                      ],
                      [
                        "to_do",
                        "$TO_DO"
                      ],
                      [
                        "fix",
                        "$FIX"
                      ],
                      [
                        "info",
                        "$INFO"
                      ],
                      [
                        "warning",
                        "$WARNING"
                      ],
                      [
                        "error",
                        "$ERROR"
                      ]
                    ]
                  },
                  {
                    "name": "LOG_LEVEL",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_debug_log_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "log"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "as"
              },
              {
                "type": "field_checkbox",
                "name": "log_level_edit"
              },
              {
                "name": "LOG_LEVEL",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "determiner",
            "output": [
              "determiner"
            ],
            "message0": "determiner",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Determiners: modify a word they are associated to designate specificity or, sometimes, a count.  For instance: \"some\" fish hooks, \"a\" pineapple, \"75\" triangles, \"our\" Trevor.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_determiner_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "DETERMINER",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "diff_of",
            "output": [
              "diff_of",
              "number_eval"
            ],
            "message0": "dec",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Subtract two numbers.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_diff_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "a"
                  },
                  {
                    "name": "A",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by"
                  },
                  {
                    "name": "B",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_diff_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "dec"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "by"
              },
              {
                "type": "field_checkbox",
                "name": "b_edit"
              },
              {
                "name": "B",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "directive",
            "output": [
              "directive",
              "grammar_maker"
            ],
            "message0": "directive",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "starts a parser scanner.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_directive_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "lede"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LEDE",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "scans"
                  },
                  {
                    "name": "SCANS",
                    "type": "input_value",
                    "checks": [
                      "scanner_maker"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_directive_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "directive"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "lede"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "lede_edit"
              },
              {
                "name": "LEDE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "scans"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "scans_edit"
              },
              {
                "name": "SCANS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "_do_nothing_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "do_nothing",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Statement which does nothing.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_do_nothing_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "why"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "REASON",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "do_nothing",
            "output": [
              "do_nothing"
            ],
            "message0": "do_nothing",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Statement which does nothing.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_do_nothing_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "why"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "REASON",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_do_nothing_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "do_nothing"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "why"
              },
              {
                "type": "field_checkbox",
                "name": "reason_edit"
              },
              {
                "name": "REASON",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "during",
            "output": [
              "during",
              "bool_eval",
              "number_eval"
            ],
            "message0": "during",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Decide whether a pattern is running.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_during_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "pattern"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PATTERN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_during_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "during"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_aliases",
            "output": [
              "eph_aliases",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_aliases_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "understand"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SHORT_NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ALIASES",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_aliases_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "as"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "aliases_edit"
              },
              {
                "name": "ALIASES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "eph_always",
            "output": [
              "eph_always"
            ],
            "message0": "eph_always",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_eph_always_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "always",
                        "$ALWAYS"
                      ]
                    ]
                  },
                  {
                    "name": "EPH_ALWAYS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "eph_aspects",
            "output": [
              "eph_aspects",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_aspects_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "aspects"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASPECTS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "traits"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TRAITS",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_aspects_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "traits"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "traits_edit"
              },
              {
                "name": "TRAITS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "eph_at",
            "output": [
              "eph_at"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_at_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "at"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "AT",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "eph"
                  },
                  {
                    "name": "EPH",
                    "type": "input_value",
                    "checks": [
                      "ephemera"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_at_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_begin_domain",
            "output": [
              "eph_begin_domain",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_begin_domain_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "domain"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "requires"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "REQUIRES",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_begin_domain_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "requires"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "requires_edit"
              },
              {
                "name": "REQUIRES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "eph_cardinality",
            "output": [
              "eph_cardinality"
            ],
            "message0": "eph_cardinality",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_eph_cardinality_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one_one",
                        "$ONE_ONE"
                      ],
                      [
                        "one_many",
                        "$ONE_MANY"
                      ],
                      [
                        "many_one",
                        "$MANY_ONE"
                      ],
                      [
                        "many_many",
                        "$MANY_MANY"
                      ]
                    ],
                    "swaps": {
                      "$ONE_ONE": "one_one",
                      "$ONE_MANY": "one_many",
                      "$MANY_ONE": "many_one",
                      "$MANY_MANY": "many_many"
                    }
                  },
                  {
                    "name": "EPH_CARDINALITY",
                    "type": "input_value",
                    "checks": [
                      "one_one",
                      "one_many",
                      "many_one",
                      "many_many"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "eph_checks",
            "output": [
              "eph_checks",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_checks_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "check"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "expect"
                  },
                  {
                    "name": "EXPECT",
                    "type": "input_value",
                    "checks": [
                      "literal_value"
                    ],
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "EXE",
                    "type": "input_statement",
                    "checks": [
                      "_execute_stack"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_checks_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "expect"
              },
              {
                "type": "field_checkbox",
                "name": "expect_edit"
              },
              {
                "name": "EXPECT",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "eph_directives",
            "output": [
              "eph_directives",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_directives_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "go"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "parse"
                  },
                  {
                    "name": "DIRECTIVE",
                    "type": "input_value",
                    "checks": [
                      "directive"
                    ],
                    "shadow": "directive"
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_directives_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_end_domain",
            "output": [
              "eph_end_domain",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_end_domain_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "domain"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_end_domain_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_kinds",
            "output": [
              "eph_kinds",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_kinds_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "FROM",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "contain"
                  },
                  {
                    "name": "CONTAIN",
                    "type": "input_value",
                    "checks": [
                      "eph_params"
                    ],
                    "shadow": "eph_params",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_kinds_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "contain"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "contain_edit"
              },
              {
                "name": "CONTAIN",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "eph_list",
            "output": [
              "eph_list"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_list_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "name": "ALL",
                    "type": "input_value",
                    "checks": [
                      "eph_at"
                    ],
                    "shadow": "eph_at",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_list_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "list"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "all_edit"
              },
              {
                "name": "ALL",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "eph_nouns",
            "output": [
              "eph_nouns",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_nouns_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "noun"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NOUN",
                    "type": "input_dummy"
                  }
                ],
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
                ]
              ]
            }
          },
          {
            "type": "_eph_nouns_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_opposites",
            "output": [
              "eph_opposites",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_opposites_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "opposite"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OPPOSITE",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "word"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "WORD",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_opposites_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_params",
            "output": [
              "eph_params"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_params_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "have"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "bool",
                        "$BOOL"
                      ],
                      [
                        "number",
                        "$NUMBER"
                      ],
                      [
                        "num_list",
                        "$NUM_LIST"
                      ],
                      [
                        "text",
                        "$TEXT"
                      ],
                      [
                        "text_list",
                        "$TEXT_LIST"
                      ],
                      [
                        "record",
                        "$RECORD"
                      ],
                      [
                        "record_list",
                        "$RECORD_LIST"
                      ]
                    ]
                  },
                  {
                    "name": "AFFINITY",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "called"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "CLASS",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "initially"
                  },
                  {
                    "name": "INITIALLY",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_params_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "class_edit"
              },
              {
                "name": "CLASS",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "initially"
              },
              {
                "type": "field_checkbox",
                "name": "initially_edit"
              },
              {
                "name": "INITIALLY",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "eph_patterns",
            "output": [
              "eph_patterns",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_patterns_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "pattern"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "with"
                  },
                  {
                    "name": "PARAMS",
                    "type": "input_value",
                    "checks": [
                      "eph_params"
                    ],
                    "shadow": "eph_params",
                    "optional": true,
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "locals"
                  },
                  {
                    "name": "LOCALS",
                    "type": "input_value",
                    "checks": [
                      "eph_params"
                    ],
                    "shadow": "eph_params",
                    "optional": true,
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "result"
                  },
                  {
                    "name": "RESULT",
                    "type": "input_value",
                    "checks": [
                      "eph_params"
                    ],
                    "shadow": "eph_params",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_patterns_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "with"
              },
              {
                "type": "field_number",
                "min": 0,
                "precision": 1,
                "name": "params_edit"
              },
              {
                "name": "PARAMS",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "locals"
              },
              {
                "type": "field_number",
                "min": 0,
                "precision": 1,
                "name": "locals_edit"
              },
              {
                "name": "LOCALS",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "result"
              },
              {
                "type": "field_checkbox",
                "name": "result_edit"
              },
              {
                "name": "RESULT",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8%9%10%11"
          },
          {
            "type": "eph_plurals",
            "output": [
              "eph_plurals",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_plurals_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "plural"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "singular"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SINGULAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_plurals_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_refs",
            "output": [
              "eph_refs",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_refs_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "refs"
                  },
                  {
                    "name": "REFS",
                    "type": "input_value",
                    "checks": [
                      "ephemera"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_refs_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "refs"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "refs_edit"
              },
              {
                "name": "REFS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "eph_relations",
            "output": [
              "eph_relations",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_relations_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "rel"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "REL",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "relate"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one_one",
                        "$ONE_ONE"
                      ],
                      [
                        "one_many",
                        "$ONE_MANY"
                      ],
                      [
                        "many_one",
                        "$MANY_ONE"
                      ],
                      [
                        "many_many",
                        "$MANY_MANY"
                      ]
                    ],
                    "swaps": {
                      "$ONE_ONE": "one_one",
                      "$ONE_MANY": "one_many",
                      "$MANY_ONE": "many_one",
                      "$MANY_MANY": "many_many"
                    }
                  },
                  {
                    "name": "CARDINALITY",
                    "type": "input_value",
                    "checks": [
                      "one_one",
                      "one_many",
                      "many_one",
                      "many_many"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_relations_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_relatives",
            "output": [
              "eph_relatives",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_relatives_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "rel"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "REL",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "relates"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NOUN",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "to"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OTHER_NOUN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_relatives_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "eph_rules",
            "output": [
              "eph_rules",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_rules_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "pattern"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "target"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TARGET",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "if"
                  },
                  {
                    "name": "FILTER",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "when"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "before",
                        "$BEFORE"
                      ],
                      [
                        "during",
                        "$DURING"
                      ],
                      [
                        "after",
                        "$AFTER"
                      ],
                      [
                        "later",
                        "$LATER"
                      ]
                    ]
                  },
                  {
                    "name": "WHEN",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "EXE",
                    "type": "input_statement",
                    "checks": [
                      "_execute_stack"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "touch"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "always",
                        "$ALWAYS"
                      ]
                    ]
                  },
                  {
                    "name": "TOUCH",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_rules_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "target"
              },
              {
                "type": "field_checkbox",
                "name": "target_edit"
              },
              {
                "name": "TARGET",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "touch"
              },
              {
                "type": "field_checkbox",
                "name": "touch_edit"
              },
              {
                "name": "TOUCH",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "eph_timing",
            "output": [
              "eph_timing"
            ],
            "message0": "eph_timing",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_eph_timing_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "before",
                        "$BEFORE"
                      ],
                      [
                        "during",
                        "$DURING"
                      ],
                      [
                        "after",
                        "$AFTER"
                      ],
                      [
                        "later",
                        "$LATER"
                      ]
                    ]
                  },
                  {
                    "name": "EPH_TIMING",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "eph_values",
            "output": [
              "eph_values",
              "ephemera"
            ],
            "message0": "eph",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_eph_values_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "noun"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NOUN",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "has"
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
                    "text": "path"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PATH",
                    "type": "input_dummy",
                    "optional": true,
                    "repeats": true
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
                    "checks": [
                      "literal_value"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_eph_values_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "eph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "path"
              },
              {
                "type": "field_number",
                "min": 0,
                "precision": 1,
                "name": "path_edit"
              },
              {
                "name": "PATH",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "equal",
            "output": [
              "equal",
              "comparator"
            ],
            "message0": "equals",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Two values exactly match.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_equal_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_equal_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "equals"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_erase_edge_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "erase",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Erase at edge: Remove one or more values from a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_erase_edge_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_front"
                  },
                  {
                    "name": "AT_EDGE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "erase_edge",
            "output": [
              "erase_edge"
            ],
            "message0": "erase",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Erase at edge: Remove one or more values from a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_erase_edge_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_front"
                  },
                  {
                    "name": "AT_EDGE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_erase_edge_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "erase"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "at_front"
              },
              {
                "type": "field_checkbox",
                "name": "at_edge_edit"
              },
              {
                "name": "AT_EDGE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_erase_index_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "erase",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Erase at index: Remove one or more values from a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_erase_index_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "count"
                  },
                  {
                    "name": "COUNT",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_index"
                  },
                  {
                    "name": "AT_INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "erase_index",
            "output": [
              "erase_index"
            ],
            "message0": "erase",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Erase at index: Remove one or more values from a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_erase_index_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "count"
                  },
                  {
                    "name": "COUNT",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_index"
                  },
                  {
                    "name": "AT_INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_erase_index_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "erase"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_erasing_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "erasing",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Erase elements from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_erasing_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "count"
                  },
                  {
                    "name": "COUNT",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_index"
                  },
                  {
                    "name": "AT_INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "AS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "erasing",
            "output": [
              "erasing"
            ],
            "message0": "erasing",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Erase elements from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_erasing_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "count"
                  },
                  {
                    "name": "COUNT",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_index"
                  },
                  {
                    "name": "AT_INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "AS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_erasing_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "erasing"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_erasing_edge_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "erasing",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Erase one element from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_erasing_edge_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_front"
                  },
                  {
                    "name": "AT_EDGE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "AS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "erasing_edge",
            "output": [
              "erasing_edge"
            ],
            "message0": "erasing",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Erase one element from the front or back of a list. Runs an activity with a list containing the erased values; the list can be empty if nothing was erased.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_erasing_edge_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_front"
                  },
                  {
                    "name": "AT_EDGE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "AS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_erasing_edge_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "erasing"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "at_front"
              },
              {
                "type": "field_checkbox",
                "name": "at_edge_edit"
              },
              {
                "name": "AT_EDGE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "else"
              },
              {
                "type": "field_checkbox",
                "name": "else_edit"
              },
              {
                "name": "ELSE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "_event_block_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "event_block",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Declare event listeners.\nListeners let objects in the game world react to changes before, during, or after they happen.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_event_block_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "the target"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "kinds",
                        "$KINDS"
                      ],
                      [
                        "noun",
                        "$NOUN"
                      ]
                    ],
                    "swaps": {
                      "$KINDS": "plural_kinds",
                      "$NOUN": "named_noun"
                    }
                  },
                  {
                    "name": "TARGET",
                    "type": "input_value",
                    "checks": [
                      "plural_kinds",
                      "named_noun"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "handlers"
                  },
                  {
                    "name": "HANDLERS",
                    "type": "input_value",
                    "checks": [
                      "event_handler"
                    ],
                    "shadow": "event_handler",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "event_block",
            "output": [
              "event_block"
            ],
            "message0": "event_block",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Declare event listeners.\nListeners let objects in the game world react to changes before, during, or after they happen.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_event_block_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "the target"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "kinds",
                        "$KINDS"
                      ],
                      [
                        "noun",
                        "$NOUN"
                      ]
                    ],
                    "swaps": {
                      "$KINDS": "plural_kinds",
                      "$NOUN": "named_noun"
                    }
                  },
                  {
                    "name": "TARGET",
                    "type": "input_value",
                    "checks": [
                      "plural_kinds",
                      "named_noun"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "handlers"
                  },
                  {
                    "name": "HANDLERS",
                    "type": "input_value",
                    "checks": [
                      "event_handler"
                    ],
                    "shadow": "event_handler",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_event_block_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "event_block"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "handlers"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "handlers_edit"
              },
              {
                "name": "HANDLERS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "event_handler",
            "output": [
              "event_handler"
            ],
            "message0": "event_handler",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_event_handler_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "event_phase"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "before",
                        "$BEFORE"
                      ],
                      [
                        "during",
                        "$WHILE"
                      ],
                      [
                        "after",
                        "$AFTER"
                      ]
                    ]
                  },
                  {
                    "name": "EVENT_PHASE",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "the event"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "EVENT",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "provides"
                  },
                  {
                    "name": "PROVIDES",
                    "type": "input_value",
                    "checks": [
                      "pattern_locals"
                    ],
                    "shadow": "pattern_locals",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "pattern_rules"
                  },
                  {
                    "name": "PATTERN_RULES",
                    "type": "input_value",
                    "checks": [
                      "pattern_rules"
                    ],
                    "shadow": "pattern_rules"
                  }
                ]
              ]
            }
          },
          {
            "type": "_event_handler_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "event_handler"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "provides"
              },
              {
                "type": "field_checkbox",
                "name": "provides_edit"
              },
              {
                "name": "PROVIDES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "event_name",
            "output": [
              "event_name"
            ],
            "message0": "event_name",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_event_name_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "EVENT_NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "event_phase",
            "output": [
              "event_phase"
            ],
            "message0": "event_phase",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_event_phase_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "before",
                        "$BEFORE"
                      ],
                      [
                        "during",
                        "$WHILE"
                      ],
                      [
                        "after",
                        "$AFTER"
                      ]
                    ]
                  },
                  {
                    "name": "EVENT_PHASE",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "event_target",
            "output": [
              "event_target"
            ],
            "message0": "event_target",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_event_target_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "kinds",
                        "$KINDS"
                      ],
                      [
                        "noun",
                        "$NOUN"
                      ]
                    ],
                    "swaps": {
                      "$KINDS": "plural_kinds",
                      "$NOUN": "named_noun"
                    }
                  },
                  {
                    "name": "EVENT_TARGET",
                    "type": "input_value",
                    "checks": [
                      "plural_kinds",
                      "named_noun"
                    ]
                  }
                ]
              ]
            }
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
              "mui": "_field_value_mutator",
              "shapeDef": [
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
                    "checks": [
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
              "mui": "_field_values_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "contains"
                  },
                  {
                    "name": "CONTAINS",
                    "type": "input_value",
                    "checks": [
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
            "type": "flow_spec",
            "output": [
              "flow_spec"
            ],
            "message0": "flow",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_flow_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "phrase"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PHRASE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trim"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
                  },
                  {
                    "name": "TRIM",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "uses"
                  },
                  {
                    "name": "TERMS",
                    "type": "input_value",
                    "checks": [
                      "term_spec"
                    ],
                    "shadow": "term_spec",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_flow_spec_mutator",
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
                "text": "name"
              },
              {
                "type": "field_checkbox",
                "name": "name_edit"
              },
              {
                "name": "NAME",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "phrase"
              },
              {
                "type": "field_checkbox",
                "name": "phrase_edit"
              },
              {
                "name": "PHRASE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "trim"
              },
              {
                "type": "field_checkbox",
                "name": "trim_edit"
              },
              {
                "name": "TRIM",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "uses"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "terms_edit"
              },
              {
                "name": "TERMS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8%9%10%11%12%13%14"
          },
          {
            "type": "from_bool",
            "output": [
              "from_bool",
              "assignment"
            ],
            "message0": "from_bool",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Assigns the calculated boolean value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_bool_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "val"
                  },
                  {
                    "name": "VAL",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_bool_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "from_bool"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_num",
            "output": [
              "from_num",
              "assignment"
            ],
            "message0": "from_num",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Assigns the calculated number.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_num_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "val"
                  },
                  {
                    "name": "VAL",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_num_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "from_num"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_num_list",
            "output": [
              "from_num_list",
              "list_source"
            ],
            "message0": "var_of_nums",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Uses a list of numbers.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_num_list_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_num_list_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "var_of_nums"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_numbers",
            "output": [
              "from_numbers",
              "assignment"
            ],
            "message0": "from_nums",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Assigns the calculated numbers.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_numbers_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "vals"
                  },
                  {
                    "name": "VALS",
                    "type": "input_value",
                    "checks": [
                      "num_list_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_numbers_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "from_nums"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_obj",
            "output": [
              "from_obj",
              "from_source_fields"
            ],
            "message0": "obj_fields",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Targets an object with a computed name.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_obj_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_obj_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "obj_fields"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_rec",
            "output": [
              "from_rec",
              "from_source_fields"
            ],
            "message0": "rec_fields",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Targets a record stored in a record.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_rec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "rec"
                  },
                  {
                    "name": "REC",
                    "type": "input_value",
                    "checks": [
                      "record_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_rec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "rec_fields"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_rec_list",
            "output": [
              "from_rec_list",
              "list_source"
            ],
            "message0": "var_of_recs",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Uses a list of records.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_rec_list_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_rec_list_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "var_of_recs"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_record",
            "output": [
              "from_record",
              "assignment"
            ],
            "message0": "from_rec",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Assigns the calculated record.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_record_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "val"
                  },
                  {
                    "name": "VAL",
                    "type": "input_value",
                    "checks": [
                      "record_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_record_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "from_rec"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_records",
            "output": [
              "from_records",
              "assignment"
            ],
            "message0": "from_recs",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Assigns the calculated records.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_records_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "vals"
                  },
                  {
                    "name": "VALS",
                    "type": "input_value",
                    "checks": [
                      "record_list_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_records_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "from_recs"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_text",
            "output": [
              "from_text",
              "assignment"
            ],
            "message0": "from_txt",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Assigns the calculated piece of text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "val"
                  },
                  {
                    "name": "VAL",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "from_txt"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_texts",
            "output": [
              "from_texts",
              "assignment"
            ],
            "message0": "from_txts",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Assigns the calculated texts.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_texts_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "vals"
                  },
                  {
                    "name": "VALS",
                    "type": "input_value",
                    "checks": [
                      "text_list_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_texts_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "from_txts"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_txt_list",
            "output": [
              "from_txt_list",
              "list_source"
            ],
            "message0": "var_of_txts",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Uses a list of text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_txt_list_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_txt_list_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "var_of_txts"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "from_var",
            "output": [
              "from_var",
              "from_source_fields"
            ],
            "message0": "var_fields",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Targets a record stored in a variable.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_from_var_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_from_var_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "var_fields"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "get_at_field",
            "output": [
              "get_at_field",
              "assignment",
              "bool_eval",
              "number_eval",
              "text_eval",
              "record_eval",
              "num_list_eval",
              "text_list_eval",
              "record_list_eval"
            ],
            "message0": "get",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Get a value from a record.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_get_at_field_mutator",
              "shapeDef": [
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
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "from_source_fields"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_get_at_field_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "get"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "get_var",
            "output": [
              "get_var",
              "assignment",
              "bool_eval",
              "number_eval",
              "text_eval",
              "record_eval",
              "num_list_eval",
              "text_list_eval",
              "record_list_eval"
            ],
            "message0": "var",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "Get Variable: Return the value of the named variable.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_get_var_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_get_var_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "var"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "grammar",
            "output": [
              "grammar"
            ],
            "message0": "grammar",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Read what the player types and turn it into actions.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_grammar_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "grammar"
                  },
                  {
                    "name": "GRAMMAR",
                    "type": "input_value",
                    "checks": [
                      "grammar_maker"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_grammar_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "grammar"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "grammar_decl",
            "output": [
              "grammar_decl"
            ],
            "message0": "grammar_decl",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Read what the player types and turn it into actions.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_grammar_decl_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "grammar"
                  },
                  {
                    "name": "GRAMMAR",
                    "type": "input_value",
                    "checks": [
                      "grammar_maker"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_grammar_decl_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "grammar_decl"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "greater_than",
            "output": [
              "greater_than",
              "comparator"
            ],
            "message0": "greater_than",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The first value is larger than the second value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_greater_than_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_greater_than_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "greater_than"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "group_spec",
            "output": [
              "group_spec"
            ],
            "message0": "group",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_group_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "contains"
                  },
                  {
                    "name": "SPECS",
                    "type": "input_value",
                    "checks": [
                      "type_spec"
                    ],
                    "shadow": "type_spec",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_group_spec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "group"
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
                "name": "specs_edit"
              },
              {
                "name": "SPECS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "has_dominion",
            "output": [
              "has_dominion",
              "bool_eval"
            ],
            "message0": "has_dominion",
            "colour": "%{BKY_LOGIC_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_has_dominion_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_has_dominion_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "has_dominion"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "has_trait",
            "output": [
              "has_trait",
              "bool_eval"
            ],
            "message0": "get",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Return true if the object is currently in the requested state.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_has_trait_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "obj"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait"
                  },
                  {
                    "name": "TRAIT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_has_trait_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "get"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "id_of",
            "output": [
              "id_of",
              "text_eval"
            ],
            "message0": "id_of",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "A unique object identifier.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_id_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_id_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "id_of"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "includes",
            "output": [
              "includes",
              "bool_eval"
            ],
            "message0": "contains",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "True if text contains text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_includes_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "part"
                  },
                  {
                    "name": "PART",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_includes_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "contains"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "into_num_list",
            "output": [
              "into_num_list",
              "list_target"
            ],
            "message0": "into_nums",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Targets a list of numbers.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_into_num_list_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_into_num_list_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "into_nums"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "into_obj",
            "output": [
              "into_obj",
              "into_target_fields"
            ],
            "message0": "obj_field",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Targets an object with a computed name.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_into_obj_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_into_obj_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "obj_field"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "into_rec_list",
            "output": [
              "into_rec_list",
              "list_target"
            ],
            "message0": "into_recs",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Targets a list of records.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_into_rec_list_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_into_rec_list_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "into_recs"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "into_txt_list",
            "output": [
              "into_txt_list",
              "list_target"
            ],
            "message0": "into_txts",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Targets a list of text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_into_txt_list_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_into_txt_list_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "into_txts"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "into_var",
            "output": [
              "into_var",
              "into_target_fields"
            ],
            "message0": "var_field",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Targets an object or record stored in a variable.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_into_var_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_into_var_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "var_field"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "is_empty",
            "output": [
              "is_empty",
              "bool_eval"
            ],
            "message0": "is",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "True if the text is empty.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_is_empty_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "empty"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_is_empty_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "is"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "is_exact_kind_of",
            "output": [
              "is_exact_kind_of",
              "bool_eval"
            ],
            "message0": "kind_of",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "True if the object is exactly the named kind.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_is_exact_kind_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "is_exactly"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KIND",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_is_exact_kind_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kind_of"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "is_kind_of",
            "output": [
              "is_kind_of",
              "bool_eval"
            ],
            "message0": "kind_of",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "True if the object is compatible with the named kind.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_is_kind_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "is"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KIND",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_is_kind_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kind_of"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "join",
            "output": [
              "join",
              "text_eval"
            ],
            "message0": "join",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns multiple pieces of text as a single new piece of text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_join_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "sep"
                  },
                  {
                    "name": "SEP",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "parts"
                  },
                  {
                    "name": "PARTS",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_join_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "join"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "sep"
              },
              {
                "type": "field_checkbox",
                "name": "sep_edit"
              },
              {
                "name": "SEP",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "parts"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "parts_edit"
              },
              {
                "name": "PARTS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "kind_of",
            "output": [
              "kind_of",
              "text_eval"
            ],
            "message0": "kind_of",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Friendly name of the object's kind.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kind_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_kind_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kind_of"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "kind_of_noun",
            "output": [
              "kind_of_noun",
              "noun_continuation"
            ],
            "message0": "kind_of_noun",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kind_of_noun_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "are_an"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "are a",
                        "$AREA"
                      ],
                      [
                        "are an",
                        "$AREAN"
                      ],
                      [
                        "is",
                        "$IS"
                      ],
                      [
                        "is a",
                        "$ISA"
                      ],
                      [
                        "is an",
                        "$ISAN"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_AN",
                    "type": "input_dummy"
                  }
                ],
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
                ]
              ]
            }
          },
          {
            "type": "_kind_of_noun_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kind_of_noun"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_kind_of_relation_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "kind_of_relation",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kind_of_relation_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "relation"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "RELATION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "cardinality"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one_to_one",
                        "$ONE_TO_ONE"
                      ],
                      [
                        "one_to_many",
                        "$ONE_TO_MANY"
                      ],
                      [
                        "many_to_one",
                        "$MANY_TO_ONE"
                      ],
                      [
                        "many_to_many",
                        "$MANY_TO_MANY"
                      ]
                    ],
                    "swaps": {
                      "$ONE_TO_ONE": "one_to_one",
                      "$ONE_TO_MANY": "one_to_many",
                      "$MANY_TO_ONE": "many_to_one",
                      "$MANY_TO_MANY": "many_to_many"
                    }
                  },
                  {
                    "name": "CARDINALITY",
                    "type": "input_value",
                    "checks": [
                      "one_to_one",
                      "one_to_many",
                      "many_to_one",
                      "many_to_many"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "kind_of_relation",
            "output": [
              "kind_of_relation"
            ],
            "message0": "kind_of_relation",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kind_of_relation_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "relation"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "RELATION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "cardinality"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one_to_one",
                        "$ONE_TO_ONE"
                      ],
                      [
                        "one_to_many",
                        "$ONE_TO_MANY"
                      ],
                      [
                        "many_to_one",
                        "$MANY_TO_ONE"
                      ],
                      [
                        "many_to_many",
                        "$MANY_TO_MANY"
                      ]
                    ],
                    "swaps": {
                      "$ONE_TO_ONE": "one_to_one",
                      "$ONE_TO_MANY": "one_to_many",
                      "$MANY_TO_ONE": "many_to_one",
                      "$MANY_TO_MANY": "many_to_many"
                    }
                  },
                  {
                    "name": "CARDINALITY",
                    "type": "input_value",
                    "checks": [
                      "one_to_one",
                      "one_to_many",
                      "many_to_one",
                      "many_to_many"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_kind_of_relation_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kind_of_relation"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_kinds_have_properties_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "kinds_have_properties",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add properties to a kind",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kinds_have_properties_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "plural_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL_KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "props"
                  },
                  {
                    "name": "PROPS",
                    "type": "input_value",
                    "checks": [
                      "property_slot"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "kinds_have_properties",
            "output": [
              "kinds_have_properties"
            ],
            "message0": "kinds_have_properties",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add properties to a kind",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kinds_have_properties_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "plural_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL_KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "props"
                  },
                  {
                    "name": "PROPS",
                    "type": "input_value",
                    "checks": [
                      "property_slot"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_kinds_have_properties_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kinds_have_properties"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "props"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "props_edit"
              },
              {
                "name": "PROPS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "kinds_of",
            "output": [
              "kinds_of",
              "text_list_eval"
            ],
            "message0": "kinds_of",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "A list of compatible kinds.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kinds_of_mutator",
              "shapeDef": [
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
                ]
              ]
            }
          },
          {
            "type": "_kinds_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kinds_of"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_kinds_of_aspect_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "kinds_of_aspect",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kinds_of_aspect_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "aspect"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASPECT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "kinds_of_aspect",
            "output": [
              "kinds_of_aspect"
            ],
            "message0": "kinds_of_aspect",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kinds_of_aspect_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "aspect"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "ASPECT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_kinds_of_aspect_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kinds_of_aspect"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_kinds_of_kind_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "kinds_of_kind",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Declare a kind",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kinds_of_kind_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "plural_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL_KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "singular_kind"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SINGULAR_KIND",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "kinds_of_kind",
            "output": [
              "kinds_of_kind"
            ],
            "message0": "kinds_of_kind",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Declare a kind",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_kinds_of_kind_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "plural_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL_KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "singular_kind"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SINGULAR_KIND",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_kinds_of_kind_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kinds_of_kind"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "less_than",
            "output": [
              "less_than",
              "comparator"
            ],
            "message0": "less_than",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The first value is less than the second value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_less_than_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_less_than_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "less_than"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "lines",
            "output": [
              "lines"
            ],
            "message0": "lines",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "A sequence of characters of any length spanning multiple lines.  Paragraphs are a prime example. Generally lines are some piece of the story that will be displayed to the player.  See also: text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_lines_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LINES",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "list_at",
            "output": [
              "list_at",
              "number_eval",
              "text_eval",
              "record_eval"
            ],
            "message0": "get",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Get a value from a list. The first element is is index 1.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_at_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "name": "LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "index"
                  },
                  {
                    "name": "INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_at_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "get"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_list_each_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "repeating",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Loops over the elements in the passed list, or runs the 'else' activity if empty.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_each_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "across"
                  },
                  {
                    "name": "LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "name": "AS",
                    "type": "input_value",
                    "checks": [
                      "list_iterator"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "list_each",
            "output": [
              "list_each"
            ],
            "message0": "repeating",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Loops over the elements in the passed list, or runs the 'else' activity if empty.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_each_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "across"
                  },
                  {
                    "name": "LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "as"
                  },
                  {
                    "name": "AS",
                    "type": "input_value",
                    "checks": [
                      "list_iterator"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "else"
                  },
                  {
                    "name": "ELSE",
                    "type": "input_value",
                    "checks": [
                      "brancher"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_each_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "repeating"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "else"
              },
              {
                "type": "field_checkbox",
                "name": "else_edit"
              },
              {
                "name": "ELSE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "list_find",
            "output": [
              "list_find",
              "bool_eval",
              "number_eval"
            ],
            "message0": "find",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Search a list for a specific value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_find_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "value"
                  },
                  {
                    "name": "VALUE",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "name": "LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_find_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "find"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "list_gather",
            "output": [
              "list_gather"
            ],
            "message0": "gather",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Transform the values from a list. The named pattern gets called once for each value in the list. It get called with two parameters: 'in' as each value from the list, and 'out' as the var passed to the gather.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_gather_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "USING",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_gather_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "gather"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "list_len",
            "output": [
              "list_len",
              "number_eval"
            ],
            "message0": "len",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Determines the number of values in a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_len_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "name": "LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_len_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "len"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_list_map_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "map",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Transform the values from one list and place the results in another list. The designated pattern is called with each value from the 'from list', one value at a time.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_map_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "to_list"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TO_LIST",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from_list"
                  },
                  {
                    "name": "FROM_LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "USING_PATTERN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "list_map",
            "output": [
              "list_map"
            ],
            "message0": "map",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Transform the values from one list and place the results in another list. The designated pattern is called with each value from the 'from list', one value at a time.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_map_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "to_list"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TO_LIST",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from_list"
                  },
                  {
                    "name": "FROM_LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "USING_PATTERN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_map_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "map"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_list_reduce_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "reduce",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Transform the values from one list by combining them into a single value. The named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_reduce_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "into"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "INTO_VALUE",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from_list"
                  },
                  {
                    "name": "FROM_LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "USING_PATTERN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "list_reduce",
            "output": [
              "list_reduce"
            ],
            "message0": "reduce",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Transform the values from one list by combining them into a single value. The named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_reduce_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "into"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "INTO_VALUE",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from_list"
                  },
                  {
                    "name": "FROM_LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "USING_PATTERN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_reduce_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "reduce"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_list_reverse_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "reverse",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Reverse a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_reverse_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "name": "LIST",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "list_reverse",
            "output": [
              "list_reverse"
            ],
            "message0": "reverse",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Reverse a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_reverse_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "name": "LIST",
                    "type": "input_value",
                    "checks": [
                      "list_source"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_reverse_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "reverse"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_list_set_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "set",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Overwrite an existing value in a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_set_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LIST",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "index"
                  },
                  {
                    "name": "INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "list_set",
            "output": [
              "list_set"
            ],
            "message0": "set",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Overwrite an existing value in a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_set_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LIST",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "index"
                  },
                  {
                    "name": "INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_set_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "set"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "list_slice",
            "output": [
              "list_slice",
              "num_list_eval",
              "text_list_eval",
              "record_list_eval"
            ],
            "message0": "slice",
            "colour": "%{BKY_MATH_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_slice_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "name": "LIST",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "start"
                  },
                  {
                    "name": "START",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ],
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "end"
                  },
                  {
                    "name": "END",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_slice_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "slice"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "start"
              },
              {
                "type": "field_checkbox",
                "name": "start_edit"
              },
              {
                "name": "START",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "end"
              },
              {
                "type": "field_checkbox",
                "name": "end_edit"
              },
              {
                "name": "END",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "_list_sort_numbers_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "sort_numbers",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_sort_numbers_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by_field"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "BY_FIELD",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "descending"
                  },
                  {
                    "name": "DESCENDING",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "list_sort_numbers",
            "output": [
              "list_sort_numbers"
            ],
            "message0": "sort_numbers",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_sort_numbers_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by_field"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "BY_FIELD",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "descending"
                  },
                  {
                    "name": "DESCENDING",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_sort_numbers_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "sort_numbers"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "descending"
              },
              {
                "type": "field_checkbox",
                "name": "descending_edit"
              },
              {
                "name": "DESCENDING",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_list_sort_text_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "sort_texts",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_sort_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by_field"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "BY_FIELD",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "descending"
                  },
                  {
                    "name": "DESCENDING",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using_case"
                  },
                  {
                    "name": "USING_CASE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "list_sort_text",
            "output": [
              "list_sort_text"
            ],
            "message0": "sort_texts",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_sort_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by_field"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "BY_FIELD",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "descending"
                  },
                  {
                    "name": "DESCENDING",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using_case"
                  },
                  {
                    "name": "USING_CASE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_sort_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "sort_texts"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "descending"
              },
              {
                "type": "field_checkbox",
                "name": "descending_edit"
              },
              {
                "name": "DESCENDING",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "using_case"
              },
              {
                "type": "field_checkbox",
                "name": "using_case_edit"
              },
              {
                "name": "USING_CASE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "_list_sort_using_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "sort",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_sort_using_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "USING",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "list_sort_using",
            "output": [
              "list_sort_using"
            ],
            "message0": "sort",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_sort_using_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "var"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "using"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "USING",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_sort_using_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "sort"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_list_splice_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "splice",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Modify a list by adding and removing elements. Note: the type of the elements being added must match the type of the list. Text cant be added to a list of numbers, numbers cant be added to a list of text. If the starting index is negative, it will begin that many elements from the end of the array. If list's length + the start is less than 0, it will begin from index 0. If the remove count is missing, it removes all elements from the start to the end; if it is 0 or negative, no elements are removed.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_splice_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LIST",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "start"
                  },
                  {
                    "name": "START",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "remove"
                  },
                  {
                    "name": "REMOVE",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "insert"
                  },
                  {
                    "name": "INSERT",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "list_splice",
            "output": [
              "list_splice",
              "num_list_eval",
              "text_list_eval",
              "record_list_eval"
            ],
            "message0": "splice",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Modify a list by adding and removing elements. Note: the type of the elements being added must match the type of the list. Text cant be added to a list of numbers, numbers cant be added to a list of text. If the starting index is negative, it will begin that many elements from the end of the array. If list's length + the start is less than 0, it will begin from index 0. If the remove count is missing, it removes all elements from the start to the end; if it is 0 or negative, no elements are removed.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_list_splice_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "list"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LIST",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "start"
                  },
                  {
                    "name": "START",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "remove"
                  },
                  {
                    "name": "REMOVE",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "insert"
                  },
                  {
                    "name": "INSERT",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_list_splice_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "splice"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "logging_level",
            "output": [
              "logging_level"
            ],
            "message0": "logging_level",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_logging_level_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "note",
                        "$NOTE"
                      ],
                      [
                        "to_do",
                        "$TO_DO"
                      ],
                      [
                        "fix",
                        "$FIX"
                      ],
                      [
                        "info",
                        "$INFO"
                      ],
                      [
                        "warning",
                        "$WARNING"
                      ],
                      [
                        "error",
                        "$ERROR"
                      ]
                    ]
                  },
                  {
                    "name": "LOGGING_LEVEL",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "make_lowercase",
            "output": [
              "make_lowercase",
              "text_eval"
            ],
            "message0": "lower",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns new text, with every letter turned into lowercase. For example, 'shout' from 'SHOUT'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_lowercase_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_make_lowercase_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "lower"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_make_opposite_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "make_opposite",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The opposite of east is west.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_opposite_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "word"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "WORD",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "opposite"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OPPOSITE",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "make_opposite",
            "output": [
              "make_opposite"
            ],
            "message0": "make_opposite",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The opposite of east is west.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_opposite_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "word"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "WORD",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "opposite"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OPPOSITE",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_make_opposite_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "make_opposite"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_make_plural_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "make_plural",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The plural of person is people.\nThe plural of person is persons.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_plural_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "singular"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SINGULAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "plural"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "make_plural",
            "output": [
              "make_plural"
            ],
            "message0": "make_plural",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The plural of person is people.\nThe plural of person is persons.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_plural_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "singular"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SINGULAR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "plural"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_make_plural_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "make_plural"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "make_reversed",
            "output": [
              "make_reversed",
              "text_eval"
            ],
            "message0": "reverse",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns new text flipped back to front. For example, 'elppA' from 'Apple', or 'noon' from 'noon'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_reversed_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_make_reversed_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "reverse"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "make_sentence_case",
            "output": [
              "make_sentence_case",
              "text_eval"
            ],
            "message0": "sentence",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns new text, start each sentence with a capital letter. For example, 'Empire Apple.' from 'Empire apple.'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_sentence_case_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_make_sentence_case_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "sentence"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "make_title_case",
            "output": [
              "make_title_case",
              "text_eval"
            ],
            "message0": "title",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns new text, starting each word with a capital letter. For example, 'Empire Apple' from 'empire apple'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_title_case_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_make_title_case_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "title"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "make_uppercase",
            "output": [
              "make_uppercase",
              "text_eval"
            ],
            "message0": "upper",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns new text, with every letter turned into uppercase. For example, 'APPLE' from 'apple'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_make_uppercase_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_make_uppercase_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "upper"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "many_many",
            "output": [
              "many_many"
            ],
            "message0": "kinds",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_many_many_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "to_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OTHER_KINDS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_many_many_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kinds"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "many_one",
            "output": [
              "many_one"
            ],
            "message0": "kinds",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_many_one_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "to_kind"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OTHER_KIND",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_many_one_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kinds"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "many_to_many",
            "output": [
              "many_to_many"
            ],
            "message0": "many_to_many",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_many_to_many_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KINDS",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "other_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OTHER_KINDS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_many_to_many_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "many_to_many"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "many_to_one",
            "output": [
              "many_to_one"
            ],
            "message0": "many_to_one",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_many_to_one_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KINDS",
                    "type": "input_dummy"
                  }
                ],
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
                ]
              ]
            }
          },
          {
            "type": "_many_to_one_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "many_to_one"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "map_connection",
            "output": [
              "map_connection"
            ],
            "message0": "map_connection",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Chooses between a one-way and a two-way connection between rooms.  Generally, this only makes sense for map headings, but it at least causes 'departing' to check that a reverse connection exists.  Note: moving from one room leads you into another somewhat generically.  Sometimes its useful to position the player on entry to a new room based on where they came from.  Using, a previous room or last used door can do the trick.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_map_connection_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "arriving_at",
                        "$ARRIVING_AT"
                      ],
                      [
                        "connecting_to",
                        "$CONNECTING_TO"
                      ]
                    ]
                  },
                  {
                    "name": "MAP_CONNECTION",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_map_departing_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "departing",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Leaving a room by by going through a door ( ex. departing the house via the front door... ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_map_departing_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "ROOM",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "name": "DOOR",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "arriving_at",
                        "$ARRIVING_AT"
                      ],
                      [
                        "connecting_to",
                        "$CONNECTING_TO"
                      ]
                    ]
                  },
                  {
                    "name": "MAP_CONNECTION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "other_room"
                  },
                  {
                    "name": "OTHER_ROOM",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ]
              ]
            }
          },
          {
            "type": "map_departing",
            "output": [
              "map_departing"
            ],
            "message0": "departing",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Leaving a room by by going through a door ( ex. departing the house via the front door... ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_map_departing_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "ROOM",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "name": "DOOR",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "arriving_at",
                        "$ARRIVING_AT"
                      ],
                      [
                        "connecting_to",
                        "$CONNECTING_TO"
                      ]
                    ]
                  },
                  {
                    "name": "MAP_CONNECTION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "other_room"
                  },
                  {
                    "name": "OTHER_ROOM",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ]
              ]
            }
          },
          {
            "type": "_map_departing_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "departing"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "map_direction",
            "output": [
              "map_direction"
            ],
            "message0": "map_direction",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "A heading for movement within the game, often connecting one room within the game to another.  The most commonly used are standard compass directions like 'north', 'east', 'south', and 'west'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_map_direction_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "MAP_DIRECTION",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_map_heading_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "heading",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Leaving a room by moving in a compass direction ( ex. heading east... ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_map_heading_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "dir"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "DIR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "ROOM",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "name": "DOOR",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "arriving_at",
                        "$ARRIVING_AT"
                      ],
                      [
                        "connecting_to",
                        "$CONNECTING_TO"
                      ]
                    ]
                  },
                  {
                    "name": "MAP_CONNECTION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "other_room"
                  },
                  {
                    "name": "OTHER_ROOM",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ]
              ]
            }
          },
          {
            "type": "map_heading",
            "output": [
              "map_heading"
            ],
            "message0": "heading",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Leaving a room by moving in a compass direction ( ex. heading east... ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_map_heading_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "dir"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "DIR",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "ROOM",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "name": "DOOR",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "arriving_at",
                        "$ARRIVING_AT"
                      ],
                      [
                        "connecting_to",
                        "$CONNECTING_TO"
                      ]
                    ]
                  },
                  {
                    "name": "MAP_CONNECTION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "other_room"
                  },
                  {
                    "name": "OTHER_ROOM",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun"
                  }
                ]
              ]
            }
          },
          {
            "type": "_map_heading_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "heading"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "via"
              },
              {
                "type": "field_checkbox",
                "name": "door_edit"
              },
              {
                "name": "DOOR",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "matches",
            "output": [
              "matches",
              "bool_eval"
            ],
            "message0": "matches",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Determine whether the specified text is similar to the specified regular expression.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_matches_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "to"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PATTERN",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_matches_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "matches"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "name_of",
            "output": [
              "name_of",
              "text_eval"
            ],
            "message0": "name_of",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Full name of the object.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_name_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_name_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "name_of"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "named_noun",
            "output": [
              "named_noun"
            ],
            "message0": "named_noun",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_named_noun_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "determiner"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "DETERMINER",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_named_noun_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "named_noun"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "named_property",
            "output": [
              "named_property"
            ],
            "message0": "named_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_named_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_named_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "named_property"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "never",
            "output": [
              "never",
              "bool_eval"
            ],
            "message0": "never",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Returns false.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_never_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_never_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "never"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_newline_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "br",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Start a new line.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_newline_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "newline",
            "output": [
              "newline"
            ],
            "message0": "br",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Start a new line.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_newline_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_newline_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "br"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_next_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "next",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "In a repeating loop, try the next iteration of the loop.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_next_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "next",
            "output": [
              "next"
            ],
            "message0": "next",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "In a repeating loop, try the next iteration of the loop.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_next_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_next_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "next"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "not",
            "output": [
              "not",
              "bool_eval"
            ],
            "message0": "not",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Returns the opposite value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_not_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test"
                  },
                  {
                    "name": "TEST",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_not_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "not"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "noun",
            "output": [
              "noun",
              "scanner_maker"
            ],
            "message0": "noun",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "makes a parser scanner.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_mutator",
              "shapeDef": [
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
                ]
              ]
            }
          },
          {
            "type": "_noun_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "noun"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_noun_assignment_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "noun_assignment",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Assign text to a noun.\nAssign text.\nGives a noun one or more lines of text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_assignment_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "property"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PROPERTY",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "the text"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LINES",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "noun_assignment",
            "output": [
              "noun_assignment"
            ],
            "message0": "noun_assignment",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Assign text to a noun.\nAssign text.\nGives a noun one or more lines of text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_assignment_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "property"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PROPERTY",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "the text"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LINES",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_noun_assignment_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "noun_assignment"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "nouns"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "nouns_edit"
              },
              {
                "name": "NOUNS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_noun_kind_statement_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "noun_kind_statement",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_kind_statement_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "is a kind"
                  },
                  {
                    "name": "KIND_OF_NOUN",
                    "type": "input_value",
                    "checks": [
                      "kind_of_noun"
                    ],
                    "shadow": "kind_of_noun"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "MORE",
                    "type": "input_value",
                    "checks": [
                      "noun_continuation"
                    ],
                    "optional": true,
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "noun_kind_statement",
            "output": [
              "noun_kind_statement"
            ],
            "message0": "noun_kind_statement",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_kind_statement_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "is a kind"
                  },
                  {
                    "name": "KIND_OF_NOUN",
                    "type": "input_value",
                    "checks": [
                      "kind_of_noun"
                    ],
                    "shadow": "kind_of_noun"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "MORE",
                    "type": "input_value",
                    "checks": [
                      "noun_continuation"
                    ],
                    "optional": true,
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_noun_kind_statement_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "noun_kind_statement"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "nouns"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "nouns_edit"
              },
              {
                "name": "NOUNS",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "and"
              },
              {
                "type": "field_number",
                "min": 0,
                "precision": 1,
                "name": "more_edit"
              },
              {
                "name": "MORE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "noun_name",
            "output": [
              "noun_name"
            ],
            "message0": "noun_name",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Noun name: Some specific person, place, or thing; or, more rarely, a kind.  Proper names are usually capitalized:  For example, maybe: 'Haruki', 'Jane', or 'Toronto'.  Common names are usually not capitalized. For example, maybe: 'table', 'chair', or 'dog park'.  A set of duplicate object uses their kind. For instance: twelve 'cats'.`",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_noun_name_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NOUN_NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "noun_relation",
            "output": [
              "noun_relation",
              "noun_continuation"
            ],
            "message0": "noun_relation",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_relation_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "are_being"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "is",
                        "$IS"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_BEING",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "relation"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "RELATION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "other_nouns"
                  },
                  {
                    "name": "OTHER_NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_noun_relation_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "noun_relation"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "are_being"
              },
              {
                "type": "field_checkbox",
                "name": "are_being_edit"
              },
              {
                "name": "ARE_BEING",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "other_nouns"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "other_nouns_edit"
              },
              {
                "name": "OTHER_NOUNS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "_noun_relation_statement_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "noun_relation_statement",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_relation_statement_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "noun_relation"
                  },
                  {
                    "name": "NOUN_RELATION",
                    "type": "input_value",
                    "checks": [
                      "noun_relation"
                    ],
                    "shadow": "noun_relation"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "MORE",
                    "type": "input_value",
                    "checks": [
                      "noun_continuation"
                    ],
                    "optional": true,
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "noun_relation_statement",
            "output": [
              "noun_relation_statement"
            ],
            "message0": "noun_relation_statement",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_relation_statement_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "noun_relation"
                  },
                  {
                    "name": "NOUN_RELATION",
                    "type": "input_value",
                    "checks": [
                      "noun_relation"
                    ],
                    "shadow": "noun_relation"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "MORE",
                    "type": "input_value",
                    "checks": [
                      "noun_continuation"
                    ],
                    "optional": true,
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_noun_relation_statement_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "noun_relation_statement"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "nouns"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "nouns_edit"
              },
              {
                "name": "NOUNS",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "and"
              },
              {
                "type": "field_number",
                "min": 0,
                "precision": 1,
                "name": "more_edit"
              },
              {
                "name": "MORE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "_noun_trait_statement_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "noun_trait_statement",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_trait_statement_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "noun_traits"
                  },
                  {
                    "name": "NOUN_TRAITS",
                    "type": "input_value",
                    "checks": [
                      "noun_traits"
                    ],
                    "shadow": "noun_traits"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "MORE",
                    "type": "input_value",
                    "checks": [
                      "noun_continuation"
                    ],
                    "optional": true,
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "noun_trait_statement",
            "output": [
              "noun_trait_statement"
            ],
            "message0": "noun_trait_statement",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_trait_statement_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "noun_traits"
                  },
                  {
                    "name": "NOUN_TRAITS",
                    "type": "input_value",
                    "checks": [
                      "noun_traits"
                    ],
                    "shadow": "noun_traits"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "and"
                  },
                  {
                    "name": "MORE",
                    "type": "input_value",
                    "checks": [
                      "noun_continuation"
                    ],
                    "optional": true,
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_noun_trait_statement_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "noun_trait_statement"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "nouns"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "nouns_edit"
              },
              {
                "name": "NOUNS",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "and"
              },
              {
                "type": "field_number",
                "min": 0,
                "precision": 1,
                "name": "more_edit"
              },
              {
                "name": "MORE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "noun_traits",
            "output": [
              "noun_traits",
              "noun_continuation"
            ],
            "message0": "noun_traits",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_noun_traits_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "are_being"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "is",
                        "$IS"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_BEING",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TRAIT",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_noun_traits_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "noun_traits"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "trait"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "trait_edit"
              },
              {
                "name": "TRAIT",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "num_list_property",
            "output": [
              "num_list_property",
              "property_slot"
            ],
            "message0": "num_list_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_num_list_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "initially"
                  },
                  {
                    "name": "INITIALLY",
                    "type": "input_value",
                    "checks": [
                      "num_list_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_num_list_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "num_list_property"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "initially"
              },
              {
                "type": "field_checkbox",
                "name": "initially_edit"
              },
              {
                "name": "INITIALLY",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "num_spec",
            "output": [
              "num_spec"
            ],
            "message0": "num",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_num_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "exclusively"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
                  },
                  {
                    "name": "EXCLUSIVELY",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "uses"
                  },
                  {
                    "type": "field_number"
                  },
                  {
                    "name": "USES",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_num_spec_mutator",
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
                "text": "exclusively"
              },
              {
                "type": "field_checkbox",
                "name": "exclusively_edit"
              },
              {
                "name": "EXCLUSIVELY",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "uses"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "uses_edit"
              },
              {
                "name": "USES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
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
              "mui": "_num_value_mutator",
              "shapeDef": [
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
              "mui": "_num_values_mutator",
              "shapeDef": [
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
            "type": "number",
            "output": [
              "number"
            ],
            "message0": "number",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_number_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_number"
                  },
                  {
                    "name": "NUMBER",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "number_property",
            "output": [
              "number_property",
              "property_slot"
            ],
            "message0": "number_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_number_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "initially"
                  },
                  {
                    "name": "INITIALLY",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_number_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "number_property"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "initially"
              },
              {
                "type": "field_checkbox",
                "name": "initially_edit"
              },
              {
                "name": "INITIALLY",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "object_exists",
            "output": [
              "object_exists",
              "bool_eval"
            ],
            "message0": "is",
            "colour": "%{BKY_LOGIC_HUE}",
            "tooltip": "Returns whether there is a object of the specified name.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_object_exists_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "valid"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_object_exists_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "is"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "one_many",
            "output": [
              "one_many"
            ],
            "message0": "kind",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_one_many_mutator",
              "shapeDef": [
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
                    "text": "to_kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OTHER_KINDS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_one_many_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kind"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "one_one",
            "output": [
              "one_one"
            ],
            "message0": "kind",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_one_one_mutator",
              "shapeDef": [
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
                    "text": "to_kind"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OTHER_KIND",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_one_one_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "kind"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "one_to_many",
            "output": [
              "one_to_many"
            ],
            "message0": "one_to_many",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_one_to_many_mutator",
              "shapeDef": [
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
                    "text": "kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KINDS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_one_to_many_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "one_to_many"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "one_to_one",
            "output": [
              "one_to_one"
            ],
            "message0": "one_to_one",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_one_to_one_mutator",
              "shapeDef": [
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
                    "text": "other_kind"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OTHER_KIND",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_one_to_one_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "one_to_one"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "option_spec",
            "output": [
              "option_spec"
            ],
            "message0": "option",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_option_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "label"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LABEL",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_option_spec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "option"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "label"
              },
              {
                "type": "field_checkbox",
                "name": "label_edit"
              },
              {
                "name": "LABEL",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "paired_action",
            "output": [
              "paired_action"
            ],
            "message0": "paired_action",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_paired_action_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "kinds"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KINDS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_paired_action_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "paired_action"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "paragraph",
            "output": [
              "paragraph"
            ],
            "message0": "paragraph",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Phrases",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_paragraph_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "story_statement"
                  },
                  {
                    "name": "STORY_STATEMENT",
                    "type": "input_statement",
                    "checks": [
                      "_story_statement_stack"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_paragraph_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "paragraph"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "story_statement"
              },
              {
                "type": "field_checkbox",
                "name": "story_statement_edit"
              },
              {
                "name": "STORY_STATEMENT",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_pattern_actions_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "pattern_actions",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add actions to a pattern.\nActions to take when using a pattern.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_actions_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "pattern name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "provides"
                  },
                  {
                    "name": "PROVIDES",
                    "type": "input_value",
                    "checks": [
                      "pattern_locals"
                    ],
                    "shadow": "pattern_locals",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "pattern_rules"
                  },
                  {
                    "name": "PATTERN_RULES",
                    "type": "input_value",
                    "checks": [
                      "pattern_rules"
                    ],
                    "shadow": "pattern_rules"
                  }
                ]
              ]
            }
          },
          {
            "type": "pattern_actions",
            "output": [
              "pattern_actions"
            ],
            "message0": "pattern_actions",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add actions to a pattern.\nActions to take when using a pattern.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_actions_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "pattern name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "provides"
                  },
                  {
                    "name": "PROVIDES",
                    "type": "input_value",
                    "checks": [
                      "pattern_locals"
                    ],
                    "shadow": "pattern_locals",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "pattern_rules"
                  },
                  {
                    "name": "PATTERN_RULES",
                    "type": "input_value",
                    "checks": [
                      "pattern_rules"
                    ],
                    "shadow": "pattern_rules"
                  }
                ]
              ]
            }
          },
          {
            "type": "_pattern_actions_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "pattern_actions"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "provides"
              },
              {
                "type": "field_checkbox",
                "name": "provides_edit"
              },
              {
                "name": "PROVIDES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_pattern_decl_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "pattern_decl",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it.  Each function in a given pattern has \"guards\" which determine whether the function applies in a particular situation.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_decl_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "params"
                  },
                  {
                    "name": "PARAMS",
                    "type": "input_value",
                    "checks": [
                      "pattern_params"
                    ],
                    "shadow": "pattern_params",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "pattern_return"
                  },
                  {
                    "name": "PATTERN_RETURN",
                    "type": "input_value",
                    "checks": [
                      "pattern_return"
                    ],
                    "shadow": "pattern_return",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "pattern_decl",
            "output": [
              "pattern_decl"
            ],
            "message0": "pattern_decl",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it.  Each function in a given pattern has \"guards\" which determine whether the function applies in a particular situation.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_decl_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "params"
                  },
                  {
                    "name": "PARAMS",
                    "type": "input_value",
                    "checks": [
                      "pattern_params"
                    ],
                    "shadow": "pattern_params",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "pattern_return"
                  },
                  {
                    "name": "PATTERN_RETURN",
                    "type": "input_value",
                    "checks": [
                      "pattern_return"
                    ],
                    "shadow": "pattern_return",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_pattern_decl_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "pattern_decl"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "params"
              },
              {
                "type": "field_checkbox",
                "name": "params_edit"
              },
              {
                "name": "PARAMS",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "pattern_return"
              },
              {
                "type": "field_checkbox",
                "name": "pattern_return_edit"
              },
              {
                "name": "PATTERN_RETURN",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "pattern_flags",
            "output": [
              "pattern_flags"
            ],
            "message0": "pattern_flags",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_pattern_flags_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "continue before",
                        "$BEFORE"
                      ],
                      [
                        "continue after",
                        "$AFTER"
                      ],
                      [
                        "terminate",
                        "$TERMINATE"
                      ]
                    ]
                  },
                  {
                    "name": "PATTERN_FLAGS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "pattern_locals",
            "output": [
              "pattern_locals"
            ],
            "message0": "pattern_locals",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Storage for values used during the execution of a pattern.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_locals_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "locals"
                  },
                  {
                    "name": "LOCALS",
                    "type": "input_value",
                    "checks": [
                      "property_slot"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_pattern_locals_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "pattern_locals"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "locals"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "locals_edit"
              },
              {
                "name": "LOCALS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "pattern_name",
            "output": [
              "pattern_name"
            ],
            "message0": "pattern_name",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_pattern_name_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PATTERN_NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "pattern_params",
            "output": [
              "pattern_params"
            ],
            "message0": "pattern_params",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Pattern parameters.\nStorage for values used during the execution of a pattern.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_params_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "props"
                  },
                  {
                    "name": "PROPS",
                    "type": "input_value",
                    "checks": [
                      "property_slot"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_pattern_params_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "pattern_params"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "props"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "props_edit"
              },
              {
                "name": "PROPS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "pattern_return",
            "output": [
              "pattern_return"
            ],
            "message0": "pattern_return",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_return_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "result"
                  },
                  {
                    "name": "RESULT",
                    "type": "input_value",
                    "checks": [
                      "property_slot"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_pattern_return_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "pattern_return"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "pattern_rule",
            "output": [
              "pattern_rule"
            ],
            "message0": "pattern_rule",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Rule",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_rule_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "conditions are met"
                  },
                  {
                    "name": "GUARD",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "continue"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "continue before",
                        "$BEFORE"
                      ],
                      [
                        "continue after",
                        "$AFTER"
                      ],
                      [
                        "terminate",
                        "$TERMINATE"
                      ]
                    ]
                  },
                  {
                    "name": "FLAGS",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "actions",
                        "$ACTIVITY"
                      ]
                    ],
                    "swaps": {
                      "$ACTIVITY": "activity"
                    }
                  },
                  {
                    "name": "HOOK",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_pattern_rule_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "pattern_rule"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "continue"
              },
              {
                "type": "field_checkbox",
                "name": "flags_edit"
              },
              {
                "name": "FLAGS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "pattern_rules",
            "output": [
              "pattern_rules"
            ],
            "message0": "pattern_rules",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pattern_rules_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "pattern_rule"
                  },
                  {
                    "name": "PATTERN_RULE",
                    "type": "input_value",
                    "checks": [
                      "pattern_rule"
                    ],
                    "shadow": "pattern_rule",
                    "optional": true,
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_pattern_rules_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "pattern_rules"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "pattern_rule"
              },
              {
                "type": "field_number",
                "min": 0,
                "precision": 1,
                "name": "pattern_rule_edit"
              },
              {
                "name": "PATTERN_RULE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "pattern_type",
            "output": [
              "pattern_type"
            ],
            "message0": "pattern_type",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_pattern_type_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PATTERN_TYPE",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "plural_kinds",
            "output": [
              "plural_kinds"
            ],
            "message0": "plural_kinds",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The plural name of a type of similar nouns.  For example: animals, containers, etc.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_plural_kinds_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLURAL_KINDS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "pluralize",
            "output": [
              "pluralize",
              "text_eval"
            ],
            "message0": "plural",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns the plural form of a singular word. (ex. apples for apple. ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_pluralize_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_pluralize_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "plural"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "position",
            "output": [
              "position"
            ],
            "message0": "src",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Identifies the location of a specific command ( ex. from an .if file ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_position_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "offset"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "OFFSET",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "in"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SOURCE",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_position_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "src"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "print_num",
            "output": [
              "print_num",
              "text_eval"
            ],
            "message0": "numeral",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Writes a number using numerals, eg. '1'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_print_num_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "num"
                  },
                  {
                    "name": "NUM",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_print_num_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "numeral"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "print_num_word",
            "output": [
              "print_num_word",
              "text_eval"
            ],
            "message0": "numeral",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Writes a number in plain english: eg. 'one'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_print_num_word_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "words"
                  },
                  {
                    "name": "NUM",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_print_num_word_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "numeral"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "product_of",
            "output": [
              "product_of",
              "number_eval"
            ],
            "message0": "mul",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Multiply two numbers.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_product_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "a"
                  },
                  {
                    "name": "A",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by"
                  },
                  {
                    "name": "B",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_product_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "mul"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "program_hook",
            "output": [
              "program_hook"
            ],
            "message0": "program_hook",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_program_hook_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "actions",
                        "$ACTIVITY"
                      ]
                    ],
                    "swaps": {
                      "$ACTIVITY": "activity"
                    }
                  },
                  {
                    "name": "PROGRAM_HOOK",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "property",
            "output": [
              "property"
            ],
            "message0": "property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PROPERTY",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_put_at_field_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "put",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Put a value into the field of an record or object.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_put_at_field_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "into"
                  },
                  {
                    "name": "INTO",
                    "type": "input_value",
                    "checks": [
                      "into_target_fields"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "AT_FIELD",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "put_at_field",
            "output": [
              "put_at_field"
            ],
            "message0": "put",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Put a value into the field of an record or object.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_put_at_field_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "into"
                  },
                  {
                    "name": "INTO",
                    "type": "input_value",
                    "checks": [
                      "into_target_fields"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "AT_FIELD",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_put_at_field_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "put"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_put_edge_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "put",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add a value to a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_put_edge_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "into"
                  },
                  {
                    "name": "INTO",
                    "type": "input_value",
                    "checks": [
                      "list_target"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_front"
                  },
                  {
                    "name": "AT_EDGE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "put_edge",
            "output": [
              "put_edge"
            ],
            "message0": "put",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add a value to a list.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_put_edge_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "into"
                  },
                  {
                    "name": "INTO",
                    "type": "input_value",
                    "checks": [
                      "list_target"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_front"
                  },
                  {
                    "name": "AT_EDGE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_put_edge_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "put"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "at_front"
              },
              {
                "type": "field_checkbox",
                "name": "at_edge_edit"
              },
              {
                "name": "AT_EDGE",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_put_index_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "put",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Replace one value in a list with another.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_put_index_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "into"
                  },
                  {
                    "name": "INTO",
                    "type": "input_value",
                    "checks": [
                      "list_target"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_index"
                  },
                  {
                    "name": "AT_INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "put_index",
            "output": [
              "put_index"
            ],
            "message0": "put",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Replace one value in a list with another.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_put_index_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "assignment"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "into"
                  },
                  {
                    "name": "INTO",
                    "type": "input_value",
                    "checks": [
                      "list_target"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "at_index"
                  },
                  {
                    "name": "AT_INDEX",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_put_index_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "put"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "quotient_of",
            "output": [
              "quotient_of",
              "number_eval"
            ],
            "message0": "div",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Divide one number by another.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_quotient_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "a"
                  },
                  {
                    "name": "A",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by"
                  },
                  {
                    "name": "B",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_quotient_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "div"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "range",
            "output": [
              "range",
              "num_list_eval"
            ],
            "message0": "range",
            "colour": "%{BKY_MATH_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_range_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "to"
                  },
                  {
                    "name": "TO",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "from"
                  },
                  {
                    "name": "FROM",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ],
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by_step"
                  },
                  {
                    "name": "BY_STEP",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_range_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "range"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "from"
              },
              {
                "type": "field_checkbox",
                "name": "from_edit"
              },
              {
                "name": "FROM",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "by_step"
              },
              {
                "type": "field_checkbox",
                "name": "by_step_edit"
              },
              {
                "name": "BY_STEP",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "reciprocal_of",
            "output": [
              "reciprocal_of",
              "text_eval"
            ],
            "message0": "reciprocal",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns the implied relative of a noun (ex. the source in a one-to-many relation.).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_reciprocal_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VIA",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_reciprocal_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "reciprocal"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "reciprocals_of",
            "output": [
              "reciprocals_of",
              "text_list_eval"
            ],
            "message0": "reciprocals",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns the implied relative of a noun (ex. the sources of a many-to-many relation.).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_reciprocals_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VIA",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_reciprocals_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "reciprocals"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "record_list_property",
            "output": [
              "record_list_property",
              "property_slot"
            ],
            "message0": "record_list_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_record_list_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "initially"
                  },
                  {
                    "name": "INITIALLY",
                    "type": "input_value",
                    "checks": [
                      "record_list_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_record_list_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "record_list_property"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "initially"
              },
              {
                "type": "field_checkbox",
                "name": "initially_edit"
              },
              {
                "name": "INITIALLY",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "record_property",
            "output": [
              "record_property",
              "property_slot"
            ],
            "message0": "record_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_record_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "initially"
                  },
                  {
                    "name": "INITIALLY",
                    "type": "input_value",
                    "checks": [
                      "record_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_record_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "record_property"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "initially"
              },
              {
                "type": "field_checkbox",
                "name": "initially_edit"
              },
              {
                "name": "INITIALLY",
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
              "mui": "_record_value_mutator",
              "shapeDef": [
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
                    "checks": [
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
              "mui": "_record_values_mutator",
              "shapeDef": [
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
                    "checks": [
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
            "type": "_relate_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "relate",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Relate two nouns.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_relate_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "to"
                  },
                  {
                    "name": "TO_OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VIA",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "relate",
            "output": [
              "relate"
            ],
            "message0": "relate",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Relate two nouns.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_relate_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "to"
                  },
                  {
                    "name": "TO_OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VIA",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_relate_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "relate"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "relation_cardinality",
            "output": [
              "relation_cardinality"
            ],
            "message0": "relation_cardinality",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_relation_cardinality_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one_to_one",
                        "$ONE_TO_ONE"
                      ],
                      [
                        "one_to_many",
                        "$ONE_TO_MANY"
                      ],
                      [
                        "many_to_one",
                        "$MANY_TO_ONE"
                      ],
                      [
                        "many_to_many",
                        "$MANY_TO_MANY"
                      ]
                    ],
                    "swaps": {
                      "$ONE_TO_ONE": "one_to_one",
                      "$ONE_TO_MANY": "one_to_many",
                      "$MANY_TO_ONE": "many_to_one",
                      "$MANY_TO_MANY": "many_to_many"
                    }
                  },
                  {
                    "name": "RELATION_CARDINALITY",
                    "type": "input_value",
                    "checks": [
                      "one_to_one",
                      "one_to_many",
                      "many_to_one",
                      "many_to_many"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "relation_name",
            "output": [
              "relation_name"
            ],
            "message0": "relation_name",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_relation_name_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "RELATION_NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "relative_of",
            "output": [
              "relative_of",
              "text_eval"
            ],
            "message0": "relative",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns the relative of a noun (ex. the target of a one-to-one relation.).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_relative_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VIA",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_relative_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "relative"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_relative_to_noun_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "relative_to_noun",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Relate nouns to each other",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_relative_to_noun_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "relation"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "RELATION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "are_being"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "is",
                        "$IS"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_BEING",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "other_nouns"
                  },
                  {
                    "name": "OTHER_NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "relative_to_noun",
            "output": [
              "relative_to_noun"
            ],
            "message0": "relative_to_noun",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Relate nouns to each other",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_relative_to_noun_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "relation"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "RELATION",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "nouns"
                  },
                  {
                    "name": "NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "are_being"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "are",
                        "$ARE"
                      ],
                      [
                        "is",
                        "$IS"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_BEING",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "other_nouns"
                  },
                  {
                    "name": "OTHER_NOUNS",
                    "type": "input_value",
                    "checks": [
                      "named_noun"
                    ],
                    "shadow": "named_noun",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_relative_to_noun_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "relative_to_noun"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "nouns"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "nouns_edit"
              },
              {
                "name": "NOUNS",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "other_nouns"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "other_nouns_edit"
              },
              {
                "name": "OTHER_NOUNS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "relatives_of",
            "output": [
              "relatives_of",
              "text_list_eval"
            ],
            "message0": "relatives",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_relatives_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "via"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VIA",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "object"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_relatives_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "relatives"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "remainder_of",
            "output": [
              "remainder_of",
              "number_eval"
            ],
            "message0": "mod",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Divide one number by another, and return the remainder.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_remainder_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "a"
                  },
                  {
                    "name": "A",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by"
                  },
                  {
                    "name": "B",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_remainder_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "mod"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "render_exp",
            "output": [
              "render_exp",
              "text_eval"
            ],
            "message0": "render_exp",
            "colour": "%{BKY_TEXTS_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_render_exp_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "expression"
                  },
                  {
                    "name": "EXPRESSION",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_render_exp_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "render_exp"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "render_field",
            "output": [
              "render_field",
              "from_source_fields"
            ],
            "message0": "render_field",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_render_field_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "name": "NAME",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_render_field_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "render_field"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "render_flags",
            "output": [
              "render_flags"
            ],
            "message0": "render_flags",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_render_flags_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "as_var",
                        "$RENDER_AS_VAR"
                      ],
                      [
                        "as_obj",
                        "$RENDER_AS_OBJ"
                      ],
                      [
                        "as_any",
                        "$RENDER_AS_ANY"
                      ]
                    ]
                  },
                  {
                    "name": "RENDER_FLAGS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "render_name",
            "output": [
              "render_name",
              "text_eval"
            ],
            "message0": "render_name",
            "colour": "%{BKY_TEXTS_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_render_name_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_render_name_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "render_name"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "render_pattern",
            "output": [
              "render_pattern",
              "assignment",
              "text_eval"
            ],
            "message0": "render",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_render_pattern_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "call"
                  },
                  {
                    "name": "CALL",
                    "type": "input_value",
                    "checks": [
                      "call_pattern"
                    ],
                    "shadow": "call_pattern"
                  }
                ]
              ]
            }
          },
          {
            "type": "_render_pattern_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "render"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "render_ref",
            "output": [
              "render_ref",
              "assignment",
              "number_eval",
              "text_eval"
            ],
            "message0": "render_ref",
            "colour": "%{BKY_PROCEDURES_HUE}",
            "tooltip": "returns the value of a variable or the id of an object.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_render_ref_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "flags"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "as_var",
                        "$RENDER_AS_VAR"
                      ],
                      [
                        "as_obj",
                        "$RENDER_AS_OBJ"
                      ],
                      [
                        "as_any",
                        "$RENDER_AS_ANY"
                      ]
                    ]
                  },
                  {
                    "name": "FLAGS",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_render_ref_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "render_ref"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "render_template",
            "output": [
              "render_template",
              "text_eval"
            ],
            "message0": "render_template",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Parse text using templates.\nSee: https://github.com/ionous/iffy/wiki/Templates.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_render_template_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "template"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEMPLATE",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_render_template_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "render_template"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "response",
            "output": [
              "response",
              "text_eval"
            ],
            "message0": "response",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Generate text in a replaceable manner.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_response_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_response_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "response"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "text"
              },
              {
                "type": "field_checkbox",
                "name": "text_edit"
              },
              {
                "name": "TEXT",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "retarget",
            "output": [
              "retarget",
              "scanner_maker"
            ],
            "message0": "retarget",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "makes a parser scanner.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_retarget_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "span"
                  },
                  {
                    "name": "SPAN",
                    "type": "input_value",
                    "checks": [
                      "scanner_maker"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_retarget_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "retarget"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "span"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "span_edit"
              },
              {
                "name": "SPAN",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "reverse",
            "output": [
              "reverse",
              "scanner_maker"
            ],
            "message0": "reverse",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "makes a parser scanner.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_reverse_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "reverses"
                  },
                  {
                    "name": "REVERSES",
                    "type": "input_value",
                    "checks": [
                      "scanner_maker"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_reverse_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "reverse"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "reverses"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "reverses_edit"
              },
              {
                "name": "REVERSES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "row",
            "output": [
              "row",
              "text_eval"
            ],
            "message0": "row",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "A single line as part of a group of lines.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_row_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_row_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "row"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "rows",
            "output": [
              "rows",
              "text_eval"
            ],
            "message0": "rows",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Group text into successive lines.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_rows_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_rows_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "rows"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_say_text_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "say",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Print some bit of text to the player.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_say_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "say_text",
            "output": [
              "say_text"
            ],
            "message0": "say",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Print some bit of text to the player.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_say_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "text"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_say_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "say"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "self",
            "output": [
              "self",
              "scanner_maker"
            ],
            "message0": "self",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "makes a parser scanner which matches the player. ( the player string is just to make the composer happy. ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_self_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "player"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "PLAYER",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_self_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "self"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_set_trait_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "put",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Put an object into a particular state.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_set_trait_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "obj"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait"
                  },
                  {
                    "name": "TRAIT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "set_trait",
            "output": [
              "set_trait"
            ],
            "message0": "put",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Put an object into a particular state.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_set_trait_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "obj"
                  },
                  {
                    "name": "OBJECT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait"
                  },
                  {
                    "name": "TRAIT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_set_trait_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "put"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "shuffle_text",
            "output": [
              "shuffle_text",
              "text_eval"
            ],
            "message0": "shuffle_text",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "When called multiple times returns its inputs at random.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_shuffle_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "parts"
                  },
                  {
                    "name": "PARTS",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_shuffle_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "shuffle_text"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "parts"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "parts_edit"
              },
              {
                "name": "PARTS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "singular_kind",
            "output": [
              "singular_kind"
            ],
            "message0": "singular_kind",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Describes a type of similar nouns.  For example: an animal, a container, etc.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_singular_kind_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SINGULAR_KIND",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "singularize",
            "output": [
              "singularize",
              "text_eval"
            ],
            "message0": "singular",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Returns the singular form of a plural word. (ex. apple for apples ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_singularize_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_singularize_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "singular"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "slash_text",
            "output": [
              "slash_text",
              "text_eval"
            ],
            "message0": "slashes",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Separates words with left-leaning slashes '/'.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_slash_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_slash_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "slashes"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "slot_spec",
            "output": [
              "slot_spec"
            ],
            "message0": "slot",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_slot_spec_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_slot_spec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "slot"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_softline_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "wbr",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Start a new line ( if not already at a new line. ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_softline_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "softline",
            "output": [
              "softline"
            ],
            "message0": "wbr",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Start a new line ( if not already at a new line. ).",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_softline_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_softline_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "wbr"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "span_text",
            "output": [
              "span_text",
              "text_eval"
            ],
            "message0": "spaces",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "Writes text with spaces between words.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_span_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_span_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "spaces"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "stopping_text",
            "output": [
              "stopping_text",
              "text_eval"
            ],
            "message0": "stopping_text",
            "colour": "%{BKY_TEXTS_HUE}",
            "tooltip": "When called multiple times returns each of its inputs in turn, sticking to the last one.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_stopping_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "parts"
                  },
                  {
                    "name": "PARTS",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_stopping_text_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "stopping_text"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "parts"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "parts_edit"
              },
              {
                "name": "PARTS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "story",
            "output": [
              "story"
            ],
            "message0": "story",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_story_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "paragraph"
                  },
                  {
                    "name": "PARAGRAPH",
                    "type": "input_value",
                    "checks": [
                      "paragraph"
                    ],
                    "shadow": "paragraph",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_story_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "story"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "paragraph"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "paragraph_edit"
              },
              {
                "name": "PARAGRAPH",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "_story_break_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "story_break",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "this cheats a bit by making the signature the same as the comment marker  that allows nodes which look like comments but are actually story breaks.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_story_break_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "story_break",
            "output": [
              "story_break"
            ],
            "message0": "story_break",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "this cheats a bit by making the signature the same as the comment marker  that allows nodes which look like comments but are actually story breaks.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_story_break_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_story_break_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "story_break"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "story_lines",
            "message0": "story_lines",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_story_lines_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "lines"
                  },
                  {
                    "name": "LINES",
                    "type": "input_statement",
                    "checks": [
                      "_story_statement_stack"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_story_lines_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "story_lines"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "str_spec",
            "output": [
              "str_spec"
            ],
            "message0": "str",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_str_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "exclusively"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
                  },
                  {
                    "name": "EXCLUSIVELY",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "uses"
                  },
                  {
                    "name": "USES",
                    "type": "input_value",
                    "checks": [
                      "option_spec"
                    ],
                    "shadow": "option_spec",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_str_spec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "str"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "exclusively"
              },
              {
                "type": "field_checkbox",
                "name": "exclusively_edit"
              },
              {
                "name": "EXCLUSIVELY",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "uses"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "uses_edit"
              },
              {
                "name": "USES",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "sum_of",
            "output": [
              "sum_of",
              "number_eval"
            ],
            "message0": "inc",
            "colour": "%{BKY_MATH_HUE}",
            "tooltip": "Add two numbers.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_sum_of_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "a"
                  },
                  {
                    "name": "A",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "by"
                  },
                  {
                    "name": "B",
                    "type": "input_value",
                    "checks": [
                      "number_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_sum_of_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "inc"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "by"
              },
              {
                "type": "field_checkbox",
                "name": "b_edit"
              },
              {
                "name": "B",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "swap_spec",
            "output": [
              "swap_spec"
            ],
            "message0": "swap",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "specifies a choice between one or more other types.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_swap_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "between"
                  },
                  {
                    "name": "BETWEEN",
                    "type": "input_value",
                    "checks": [
                      "choice_spec"
                    ],
                    "shadow": "choice_spec",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_swap_spec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "swap"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "between"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "between_edit"
              },
              {
                "name": "BETWEEN",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "term_spec",
            "output": [
              "term_spec"
            ],
            "message0": "term",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_term_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "key"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "KEY",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "type"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "private"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
                  },
                  {
                    "name": "PRIVATE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "optional"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
                  },
                  {
                    "name": "OPTIONAL",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "repeats"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
                  },
                  {
                    "name": "REPEATS",
                    "type": "input_dummy",
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_term_spec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "term"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "name"
              },
              {
                "type": "field_checkbox",
                "name": "name_edit"
              },
              {
                "name": "NAME",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "type"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "private"
              },
              {
                "type": "field_checkbox",
                "name": "private_edit"
              },
              {
                "name": "PRIVATE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "optional"
              },
              {
                "type": "field_checkbox",
                "name": "optional_edit"
              },
              {
                "name": "OPTIONAL",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "repeats"
              },
              {
                "type": "field_checkbox",
                "name": "repeats_edit"
              },
              {
                "name": "REPEATS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8%9%10%11%12%13%14%15%16%17"
          },
          {
            "type": "test_bool",
            "output": [
              "test_bool"
            ],
            "message0": "test_bool",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_test_bool_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
                  },
                  {
                    "name": "TEST_BOOL",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "test_embed",
            "output": [
              "test_embed",
              "test_slot"
            ],
            "message0": "embed",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_test_embed_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test_flow"
                  },
                  {
                    "name": "TEST_FLOW",
                    "type": "input_value",
                    "checks": [
                      "test_flow"
                    ],
                    "shadow": "test_flow"
                  }
                ]
              ]
            }
          },
          {
            "type": "_test_embed_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "embed"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
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
              "mui": "_test_flow_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "slot"
                  },
                  {
                    "name": "SLOT",
                    "type": "input_value",
                    "checks": [
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
                    "type": "field_dropdown",
                    "options": [
                      [
                        "true",
                        "$TRUE"
                      ],
                      [
                        "false",
                        "$FALSE"
                      ]
                    ]
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
                        "flow",
                        "$A"
                      ],
                      [
                        "slot",
                        "$B"
                      ],
                      [
                        "text",
                        "$C"
                      ]
                    ],
                    "swaps": {
                      "$A": "test_flow",
                      "$B": "test_slot",
                      "$C": "test_txt"
                    }
                  },
                  {
                    "name": "SWAP",
                    "type": "input_value",
                    "checks": [
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
                    "checks": [
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
            "type": "test_name",
            "output": [
              "test_name"
            ],
            "message0": "test_name",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_test_name_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEST_NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "test_num",
            "output": [
              "test_num"
            ],
            "message0": "test_num",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_test_num_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_number"
                  },
                  {
                    "name": "TEST_NUM",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "test_output",
            "output": [
              "test_output",
              "testing"
            ],
            "message0": "test_output",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Expect that a test uses 'Say' to print some specific text.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_test_output_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "lines"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "LINES",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_test_output_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "test_output"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_test_rule_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "test_rule",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add actions to a test",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_test_rule_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test_name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEST_NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "actions",
                        "$ACTIVITY"
                      ]
                    ],
                    "swaps": {
                      "$ACTIVITY": "activity"
                    }
                  },
                  {
                    "name": "HOOK",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "test_rule",
            "output": [
              "test_rule"
            ],
            "message0": "test_rule",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Add actions to a test",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_test_rule_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test_name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEST_NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "actions",
                        "$ACTIVITY"
                      ]
                    ],
                    "swaps": {
                      "$ACTIVITY": "activity"
                    }
                  },
                  {
                    "name": "HOOK",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_test_rule_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "test_rule"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_test_scene_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "test_scene",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Create a scene for testing",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_test_scene_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test_name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEST_NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "story"
                  },
                  {
                    "name": "STORY",
                    "type": "input_value",
                    "checks": [
                      "story"
                    ],
                    "shadow": "story"
                  }
                ]
              ]
            }
          },
          {
            "type": "test_scene",
            "output": [
              "test_scene"
            ],
            "message0": "test_scene",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Create a scene for testing",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_test_scene_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test_name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEST_NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "story"
                  },
                  {
                    "name": "STORY",
                    "type": "input_value",
                    "checks": [
                      "story"
                    ],
                    "shadow": "story"
                  }
                ]
              ]
            }
          },
          {
            "type": "_test_scene_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "test_scene"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "_test_statement_stack",
            "nextStatement": [
              "_story_statement_stack"
            ],
            "previousStatement": [
              "_story_statement_stack"
            ],
            "message0": "test_statement",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Describe test results",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_test_statement_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test_name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEST_NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "expectation"
                  },
                  {
                    "name": "TEST",
                    "type": "input_value",
                    "checks": [
                      "testing"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "test_statement",
            "output": [
              "test_statement"
            ],
            "message0": "test_statement",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Describe test results",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_test_statement_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "test_name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEST_NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "expectation"
                  },
                  {
                    "name": "TEST",
                    "type": "input_value",
                    "checks": [
                      "testing"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "_test_statement_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "test_statement"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "test_str",
            "output": [
              "test_str"
            ],
            "message0": "test_str",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_test_str_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "one",
                        "$ONE"
                      ],
                      [
                        "other",
                        "$OTHER"
                      ],
                      [
                        "option",
                        "$OPTION"
                      ]
                    ]
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
            "type": "test_swap",
            "output": [
              "test_swap"
            ],
            "message0": "test_swap",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_test_swap_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "flow",
                        "$A"
                      ],
                      [
                        "slot",
                        "$B"
                      ],
                      [
                        "text",
                        "$C"
                      ]
                    ],
                    "swaps": {
                      "$A": "test_flow",
                      "$B": "test_slot",
                      "$C": "test_txt"
                    }
                  },
                  {
                    "name": "TEST_SWAP",
                    "type": "input_value",
                    "checks": [
                      "test_flow",
                      "test_slot",
                      "test_txt"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "test_txt",
            "output": [
              "test_txt"
            ],
            "message0": "test_txt",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_test_txt_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEST_TXT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "text",
            "output": [
              "text"
            ],
            "message0": "text",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "A sequence of characters of any length, all on one line.  Examples include letters, words, or short sentences.  Text is generally something displayed to the player.  See also: lines.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_text_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TEXT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "text_list_property",
            "output": [
              "text_list_property",
              "property_slot"
            ],
            "message0": "text_list_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_text_list_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "initially"
                  },
                  {
                    "name": "INITIALLY",
                    "type": "input_value",
                    "checks": [
                      "text_list_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_text_list_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "text_list_property"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "initially"
              },
              {
                "type": "field_checkbox",
                "name": "initially_edit"
              },
              {
                "name": "INITIALLY",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "text_property",
            "output": [
              "text_property",
              "property_slot"
            ],
            "message0": "text_property",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_text_property_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "of"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TYPE",
                    "type": "input_dummy",
                    "optional": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "initially"
                  },
                  {
                    "name": "INITIALLY",
                    "type": "input_value",
                    "checks": [
                      "text_eval"
                    ],
                    "optional": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_text_property_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "text_property"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "of"
              },
              {
                "type": "field_checkbox",
                "name": "type_edit"
              },
              {
                "name": "TYPE",
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "initially"
              },
              {
                "type": "field_checkbox",
                "name": "initially_edit"
              },
              {
                "name": "INITIALLY",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
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
              "mui": "_text_value_mutator",
              "shapeDef": [
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
              "mui": "_text_values_mutator",
              "shapeDef": [
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
          },
          {
            "type": "trait",
            "output": [
              "trait"
            ],
            "message0": "trait",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_trait_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TRAIT",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "trait_phrase",
            "output": [
              "trait_phrase"
            ],
            "message0": "trait_phrase",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_trait_phrase_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "are_either"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "can be",
                        "$CANBE"
                      ],
                      [
                        "are either",
                        "$EITHER"
                      ]
                    ]
                  },
                  {
                    "name": "ARE_EITHER",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "trait"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "TRAIT",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_trait_phrase_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "trait_phrase"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "trait"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "trait_edit"
              },
              {
                "name": "TRAIT",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          },
          {
            "type": "trigger_cycle",
            "output": [
              "trigger_cycle",
              "trigger"
            ],
            "message0": "every",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_trigger_cycle_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_trigger_cycle_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "every"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "trigger_once",
            "output": [
              "trigger_once",
              "trigger"
            ],
            "message0": "at",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_trigger_once_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_trigger_once_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "at"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "trigger_switch",
            "output": [
              "trigger_switch",
              "trigger"
            ],
            "message0": "after",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_trigger_switch_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_trigger_switch_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "after"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "type_spec",
            "output": [
              "type_spec"
            ],
            "message0": "spec",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "can optionally fit one or more slots, or be part of one or more groups.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_type_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "name"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "NAME",
                    "type": "input_dummy"
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "with"
                  },
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "flow",
                        "$FLOW"
                      ],
                      [
                        "slot",
                        "$SLOT"
                      ],
                      [
                        "swap",
                        "$SWAP"
                      ],
                      [
                        "num",
                        "$NUM"
                      ],
                      [
                        "str",
                        "$STR"
                      ],
                      [
                        "group",
                        "$GROUP"
                      ]
                    ],
                    "swaps": {
                      "$FLOW": "flow_spec",
                      "$SLOT": "slot_spec",
                      "$SWAP": "swap_spec",
                      "$NUM": "num_spec",
                      "$STR": "str_spec",
                      "$GROUP": "group_spec"
                    }
                  },
                  {
                    "name": "SPEC",
                    "type": "input_value",
                    "checks": [
                      "flow_spec",
                      "slot_spec",
                      "swap_spec",
                      "num_spec",
                      "str_spec",
                      "group_spec"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "slots"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "SLOTS",
                    "type": "input_dummy",
                    "optional": true,
                    "repeats": true
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "groups"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "GROUPS",
                    "type": "input_dummy",
                    "optional": true,
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_type_spec_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "spec"
              },
              {
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
              },
              {
                "type": "field_label",
                "text": "groups"
              },
              {
                "type": "field_number",
                "min": 0,
                "precision": 1,
                "name": "groups_edit"
              },
              {
                "name": "GROUPS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5%6%7%8"
          },
          {
            "type": "unequal",
            "output": [
              "unequal",
              "comparator"
            ],
            "message0": "other_than",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "The first value doesn't equal the second value.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_unequal_mutator",
              "shapeDef": []
            }
          },
          {
            "type": "_unequal_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "other_than"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "uses_spec",
            "output": [
              "uses_spec"
            ],
            "message0": "uses_spec",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_uses_spec_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_dropdown",
                    "options": [
                      [
                        "flow",
                        "$FLOW"
                      ],
                      [
                        "slot",
                        "$SLOT"
                      ],
                      [
                        "swap",
                        "$SWAP"
                      ],
                      [
                        "num",
                        "$NUM"
                      ],
                      [
                        "str",
                        "$STR"
                      ],
                      [
                        "group",
                        "$GROUP"
                      ]
                    ],
                    "swaps": {
                      "$FLOW": "flow_spec",
                      "$SLOT": "slot_spec",
                      "$SWAP": "swap_spec",
                      "$NUM": "num_spec",
                      "$STR": "str_spec",
                      "$GROUP": "group_spec"
                    }
                  },
                  {
                    "name": "USES_SPEC",
                    "type": "input_value",
                    "checks": [
                      "flow_spec",
                      "slot_spec",
                      "swap_spec",
                      "num_spec",
                      "str_spec",
                      "group_spec"
                    ]
                  }
                ]
              ]
            }
          },
          {
            "type": "variable_name",
            "output": [
              "variable_name"
            ],
            "message0": "variable_name",
            "colour": "%{BKY_COLOUR_HUE}",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "customData": {
              "mui": "_variable_name_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "VARIABLE_NAME",
                    "type": "input_dummy"
                  }
                ]
              ]
            }
          },
          {
            "type": "_while_stack",
            "nextStatement": [
              "_execute_stack"
            ],
            "previousStatement": [
              "_execute_stack"
            ],
            "message0": "repeating",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Keep running a series of actions while a condition is true.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_while_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "true"
                  },
                  {
                    "name": "TRUE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "while",
            "output": [
              "while"
            ],
            "message0": "repeating",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "Keep running a series of actions while a condition is true.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_while_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "true"
                  },
                  {
                    "name": "TRUE",
                    "type": "input_value",
                    "checks": [
                      "bool_eval"
                    ]
                  }
                ],
                [
                  {
                    "type": "field_label",
                    "text": "do"
                  },
                  {
                    "name": "DO",
                    "type": "input_value",
                    "checks": [
                      "activity"
                    ],
                    "shadow": "activity"
                  }
                ]
              ]
            }
          },
          {
            "type": "_while_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "repeating"
              },
              {
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2"
          },
          {
            "type": "words",
            "output": [
              "words",
              "scanner_maker"
            ],
            "message0": "words",
            "colour": "%{BKY_COLOUR_HUE}",
            "tooltip": "makes a parser scanner.",
            "extensions": [
              "tapestry_generic_mixin",
              "tapestry_generic_extension"
            ],
            "mutator": "tapestry_generic_mutation",
            "customData": {
              "mui": "_words_mutator",
              "shapeDef": [
                [
                  {
                    "type": "field_label",
                    "text": "words"
                  },
                  {
                    "type": "field_input"
                  },
                  {
                    "name": "WORDS",
                    "type": "input_dummy",
                    "repeats": true
                  }
                ]
              ]
            }
          },
          {
            "type": "_words_mutator",
            "style": "logic_blocks",
            "inputsInline": false,
            "args0": [
              {
                "type": "field_label",
                "text": "words"
              },
              {
                "type": "input_dummy"
              },
              {
                "type": "field_label",
                "text": "words"
              },
              {
                "type": "field_number",
                "min": 1,
                "precision": 1,
                "name": "words_edit"
              },
              {
                "name": "WORDS",
                "type": "input_dummy"
              }
            ],
            "message0": "%1%2%3%4%5"
          }
        ]

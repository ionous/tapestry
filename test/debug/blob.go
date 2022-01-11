package debug

var Blob = `{
  "type": "story",
  "value": {
    "$PARAGRAPH": [
      {
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "type": "story_statement",
              "value": {
                "type": "noun_kind_statement",
                "value": {
                  "$KIND_OF_NOUN": {
                    "type": "kind_of_noun",
                    "value": {
                      "$ARE_AN": {
                        "type": "are_an",
                        "value": "$ISA"
                      },
                      "$KIND": {
                        "type": "singular_kind",
                        "value": "room"
                      }
                    }
                  },
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "cabin"
                        }
                      }
                    }
                  ]
                }
              }
            }
          ]
        }
      },
      {
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "type": "story_statement",
              "value": {
                "type": "noun_relation_statement",
                "value": {
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "cabin"
                        }
                      }
                    }
                  ],
                  "$NOUN_RELATION": {
                    "type": "noun_relation",
                    "value": {
                      "$NOUNS": [
                        {
                          "type": "named_noun",
                          "value": {
                            "$DETERMINER": {
                              "type": "determiner",
                              "value": "$A"
                            },
                            "$NAME": {
                              "type": "noun_name",
                              "value": "glass case"
                            }
                          }
                        }
                      ],
                      "$RELATION": {
                        "type": "relation_name",
                        "value": "contains"
                      }
                    }
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "relative_to_noun",
                "value": {
                  "$ARE_BEING": {
                    "type": "are_being",
                    "value": "$IS"
                  },
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "glass case"
                        }
                      }
                    }
                  ],
                  "$NOUNS1": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$A"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "collection of fishing rods"
                        }
                      }
                    }
                  ],
                  "$RELATION": {
                    "type": "relation_name",
                    "value": "In"
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "noun_trait_statement",
                "value": {
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "case"
                        }
                      }
                    }
                  ],
                  "$NOUN_TRAITS": {
                    "type": "noun_traits",
                    "value": {
                      "$ARE_BEING": {
                        "type": "are_being",
                        "value": "$IS"
                      },
                      "$TRAIT": [
                        {
                          "type": "trait",
                          "value": "closed"
                        },
                        {
                          "type": "trait",
                          "value": "transparent"
                        },
                        {
                          "type": "trait",
                          "value": "lockable"
                        },
                        {
                          "type": "trait",
                          "value": "locked"
                        }
                      ]
                    }
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "noun_trait_statement",
                "value": {
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "case"
                        }
                      }
                    }
                  ],
                  "$NOUN_TRAITS": {
                    "type": "noun_traits",
                    "value": {
                      "$ARE_BEING": {
                        "type": "are_being",
                        "value": "$IS"
                      },
                      "$TRAIT": [
                        {
                          "type": "trait",
                          "value": "scenery"
                        }
                      ]
                    }
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "noun_relation_statement",
                "value": {
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "small silver key"
                        }
                      }
                    }
                  ],
                  "$NOUN_RELATION": {
                    "type": "noun_relation",
                    "value": {
                      "$NOUNS": [
                        {
                          "type": "named_noun",
                          "value": {
                            "$DETERMINER": {
                              "type": "determiner",
                              "value": "$THE"
                            },
                            "$NAME": {
                              "type": "noun_name",
                              "value": "case"
                            }
                          }
                        }
                      ],
                      "$RELATION": {
                        "type": "relation_name",
                        "value": "unlocks"
                      }
                    }
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "noun_relation_statement",
                "value": {
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "bench"
                        }
                      }
                    }
                  ],
                  "$NOUN_RELATION": {
                    "type": "noun_relation",
                    "value": {
                      "$ARE_BEING": {
                        "type": "are_being",
                        "value": "$IS"
                      },
                      "$NOUNS": [
                        {
                          "type": "named_noun",
                          "value": {
                            "$DETERMINER": {
                              "type": "determiner",
                              "value": "$THE"
                            },
                            "$NAME": {
                              "type": "noun_name",
                              "value": "cabin"
                            }
                          }
                        }
                      ],
                      "$RELATION": {
                        "type": "relation_name",
                        "value": "in"
                      }
                    }
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "relative_to_noun",
                "value": {
                  "$ARE_BEING": {
                    "type": "are_being",
                    "value": "$ARE"
                  },
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "bench"
                        }
                      }
                    }
                  ],
                  "$NOUNS1": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "some"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "blue vinyl cushions"
                        }
                      }
                    }
                  ],
                  "$RELATION": {
                    "type": "relation_name",
                    "value": "On"
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "noun_trait_statement",
                "value": {
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "bench"
                        }
                      }
                    }
                  ],
                  "$NOUN_TRAITS": {
                    "type": "noun_traits",
                    "value": {
                      "$ARE_BEING": {
                        "type": "are_being",
                        "value": "$IS"
                      },
                      "$TRAIT": [
                        {
                          "type": "trait",
                          "value": "enterable"
                        },
                        {
                          "type": "trait",
                          "value": "scenery"
                        }
                      ]
                    }
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "noun_trait_statement",
                "value": {
                  "$NOUNS": [
                    {
                      "type": "named_noun",
                      "value": {
                        "$DETERMINER": {
                          "type": "determiner",
                          "value": "$THE"
                        },
                        "$NAME": {
                          "type": "noun_name",
                          "value": "cushions"
                        }
                      }
                    }
                  ],
                  "$NOUN_TRAITS": {
                    "type": "noun_traits",
                    "value": {
                      "$ARE_BEING": {
                        "type": "are_being",
                        "value": "$ARE"
                      },
                      "$TRAIT": [
                        {
                          "type": "trait",
                          "value": "scenery"
                        }
                      ]
                    }
                  }
                }
              }
            }
          ]
        }
      }
    ]
  }
}
`

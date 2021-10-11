package jsn_test

import (
	"encoding/json"
	"hash/fnv"
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/jsn/dout"
)

func TestDetails(t *testing.T) {
	src := debug.FactorialStory
	out := dout.NewEncoder()
	src.Marshal(out)
	if d, e := out.Data(); e != nil {
		t.Fatal(e)
	} else if b, e := json.MarshalIndent(d, "", "  "); e != nil {
		t.Fatal(e)
	} else if val := hash(b); val != 0x53398df7 {
		t.Log(string(b))
		t.Fatalf("mismatched output 0x%0x", val)
	} else {
		t.Log(string(b))

		// var dst story.Story
		// dec := din.NewDecoder(makeRegistry(), []byte(det))
		// dst.Marshal(dec)
		// if _, e := dec.Data(); e != nil {
		// 	t.Log(e)
		// }
		// if diff := pretty.Diff(src, &dst); len(diff) != 0 {
		// 	pretty.Print(dst)
		// 	t.Fatal(diff)
		// }
		// }
	}
}

func makeRegistry() composer.Registry {
	var k composer.Registry
	for _, slats := range iffy.AllSlats {
		k.RegisterTypes(slats)
	}
	k.RegisterTypes(story.Slats)
	return k
}

var det = `{
	"type": "story",
	"value": {
	  "$PARAGRAPH": [
		{
		  "type": "paragraph",
		  "value": {
			"$STORY_STATEMENT": [
			  {
				"type": "test_statement",
				"value": {
				  "$TEST": {
					"type": "test_output",
					"value": {
					  "$LINES": {
						"type": "lines",
						"value": "6"
					  }
					}
				  },
				  "$TEST_NAME": {
					"type": "test_name",
					"value": "factorial"
				  }
				}
			  },
			  {
				"type": "test_rule",
				"value": {
				  "$HOOK": {
					"type": "program_hook",
					"value": {
					  "$ACTIVITY": {
						"type": "activity",
						"value": {
						  "$EXE": [
							{
							  "type": "say_text",
							  "value": {
								"$TEXT": {
								  "type": "print_num",
								  "value": {
									"$NUM": {
									  "type": "call_pattern",
									  "value": {
										"$ARGUMENTS": {
										  "type": "call_args",
										  "value": {
											"$ARGS": [
											  {
												"type": "call_arg",
												"value": {
												  "$FROM": {
													"type": "from_num",
													"value": {
													  "$VAL": {
														"type": "num_value",
														"value": {
														  "$NUM": {
															"type": "number",
															"value": 3
														  }
														}
													  }
													}
												  },
												  "$NAME": {
													"type": "text",
													"value": "num"
												  }
												}
											  }
											]
										  }
										},
										"$PATTERN": {
										  "type": "pattern_name",
										  "value": "factorial"
										}
									  }
									}
								  }
								}
							  }
							}
						  ]
						}
					  }
					}
				  },
				  "$TEST_NAME": {
					"type": "test_name",
					"value": "factorial"
				  }
				}
			  },
			  {
				"type": "pattern_decl",
				"value": {
				  "$NAME": {
					"type": "pattern_name",
					"value": "factorial"
				  },
				  "$OPTVARS": {
					"type": "pattern_variables_tail",
					"value": {
					  "$VARIABLE_DECL": [
						{
						  "type": "variable_decl",
						  "value": {
							"$AN": {
							  "type": "determiner",
							  "value": "a"
							},
							"$NAME": {
							  "type": "variable_name",
							  "value": "num"
							},
							"$TYPE": {
							  "type": "variable_type",
							  "value": {
								"$PRIMITIVE": {
								  "type": "primitive_type",
								  "value": "$NUMBER"
								}
							  }
							}
						  }
						}
					  ]
					}
				  },
				  "$TYPE": {
					"type": "pattern_type",
					"value": "$PATTERNS"
				  }
				}
			  },
			  {
				"type": "pattern_actions",
				"value": {
				  "$NAME": {
					"type": "pattern_name",
					"value": "factorial"
				  },
				  "$PATTERN_RETURN": {
					"type": "pattern_return",
					"value": {
					  "$RESULT": {
						"type": "variable_decl",
						"value": {
						  "$AN": {
							"type": "determiner",
							"value": "a"
						  },
						  "$NAME": {
							"type": "variable_name",
							"value": "num"
						  },
						  "$TYPE": {
							"type": "variable_type",
							"value": {
							  "$PRIMITIVE": {
								"type": "primitive_type",
								"value": "$NUMBER"
							  }
							}
						  }
						}
					  }
					}
				  },
				  "$PATTERN_RULES": {
					"type": "pattern_rules",
					"value": {
					  "$PATTERN_RULE": [
						{
						  "type": "pattern_rule",
						  "value": {
							"$GUARD": {
							  "type": "always",
							  "value": {}
							},
							"$HOOK": {
							  "type": "program_hook",
							  "value": {
								"$ACTIVITY": {
								  "type": "activity",
								  "value": {
									"$EXE": [
									  {
										"type": "assign",
										"value": {
										  "$FROM": {
											"type": "from_num",
											"value": {
											  "$VAL": {
												"type": "product_of",
												"value": {
												  "$A": {
													"type": "get_var",
													"value": {
													  "$NAME": {
														"type": "variable_name",
														"value": "num"
													  }
													}
												  },
												  "$B": {
													"type": "diff_of",
													"value": {
													  "$A": {
														"type": "get_var",
														"value": {
														  "$NAME": {
															"type": "variable_name",
															"value": "num"
														  }
														}
													  },
													  "$B": {
														"type": "num_value",
														"value": {
														  "$NUM": {
															"type": "number",
															"value": 1
														  }
														}
													  }
													}
												  }
												}
											  }
											}
										  },
										  "$VAR": {
											"type": "variable_name",
											"value": "num"
										  }
										}
									  }
									]
								  }
								}
							  }
							}
						  }
						}
					  ]
					}
				  }
				}
			  },
			  {
				"type": "pattern_actions",
				"value": {
				  "$NAME": {
					"type": "pattern_name",
					"value": "factorial"
				  },
				  "$PATTERN_RETURN": {
					"type": "pattern_return",
					"value": {
					  "$RESULT": {
						"type": "variable_decl",
						"value": {
						  "$AN": {
							"type": "determiner",
							"value": "a"
						  },
						  "$NAME": {
							"type": "variable_name",
							"value": "num"
						  },
						  "$TYPE": {
							"type": "variable_type",
							"value": {
							  "$PRIMITIVE": {
								"type": "primitive_type",
								"value": "$NUMBER"
							  }
							}
						  }
						}
					  }
					}
				  },
				  "$PATTERN_RULES": {
					"type": "pattern_rules",
					"value": {
					  "$PATTERN_RULE": [
						{
						  "type": "pattern_rule",
						  "value": {
							"$GUARD": {
							  "type": "compare_num",
							  "value": {
								"$A": {
								  "type": "get_var",
								  "value": {
									"$NAME": {
									  "type": "variable_name",
									  "value": "num"
									}
								  }
								},
								"$B": {
								  "type": "num_value",
								  "value": {
									"$NUM": {
									  "type": "number",
									  "value": 0
									}
								  }
								},
								"$IS": {
								  "type": "equal",
								  "value": {}
								}
							  }
							},
							"$HOOK": {
							  "type": "program_hook",
							  "value": {
								"$ACTIVITY": {
								  "type": "activity",
								  "value": {
									"$EXE": [
									  {
										"type": "assign",
										"value": {
										  "$FROM": {
											"type": "from_num",
											"value": {
											  "$VAL": {
												"type": "num_value",
												"value": {
												  "$NUM": {
													"type": "number",
													"value": 1
												  }
												}
											  }
											}
										  },
										  "$VAR": {
											"type": "variable_name",
											"value": "num"
										  }
										}
									  }
									]
								  }
								}
							  }
							}
						  }
						}
					  ]
					}
				  }
				}
			  }
			]
		  }
		}
	  ]
	}
  }`

// func TestCompact(t *testing.T) {
// 	src := debug.FactorialStory
// 	out := cout.NewEncoder()
// 	src.Marshal(out)
// 	if d, e := out.Data(); e != nil {
// 		t.Fatal(e)
// 	} else if b, e := json.MarshalIndent(d, "", "  "); e != nil {
// 		t.Fatal(e)
// 	} else if val := hash(b); val != 0xd86f0fd9 {
// 		t.Log(string(b))
// 		t.Fatalf("mismatched output 0x%0x", val)
// 	} else {

// 	}
// }

func hash(b []byte) uint32 {
	hash := fnv.New32a()
	hash.Write(b)
	return hash.Sum32()
}

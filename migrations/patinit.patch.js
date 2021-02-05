[
  {
    "op": "test",
    "path": "$..[?(@.type=='pattern_locals')]",
    "reason": "select the nodes containing type #choose.",
    "subpatches": [
      {
        "op": "move",
        "from": "$.value['$VARIABLE_DECL']",
        "path": "$.value['$LOCAL_DECL']",
        "reason": "need to move it aside to create our new node"
      },
      {
        "op": "test",
        "path": "$..[?(@.type=='variable_decl')]",
        "reason": "select the nodes containing type #choose.",
        "subpatches": [
          {
            "op": "move",
            "from": "$.value",
            "path": "$.temp",
            "reason": "need to move it aside to create our new node"
          },
          {
            "op": "add",
            "path": "$.value",
            "reason": "add the new structure.",
            "value": {
              "$VARIABLE_DECL": {
                "type": "variable_decl",
                "value": null
              }
            }
          },
          {
            "op": "move",
            "from": "$.temp",
            "path": "$.value['$VARIABLE_DECL'].value",
            "reason": "copy our saved value to our added null value"
          },
          {
            "op": "replace",
            "path": "$.type",
            "value": "local_decl",
            "reason": "rename to the new node type"
          }
        ]
      }
    ]
  }
]

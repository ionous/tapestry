[
  {
    "op": "test",
    "path": {
      "parent": "$..[?(@.type=='trying')]"
    },
    "reason": "select the nodes containing type #trying.",
    "subpatches": [
      {
        "op": "move",
        "from": "$.value['$AS']",
        "path": "$.value['$ASSIGN']",
        "reason": "rename as to assign."
      },
      {
        "op": "add",
        "path": "$.value['$FROM']",
        "value": {
          "type": "assignment",
          "value": {
            "type": "from_text",
            "value": {
              "$VAL": {
                "type": "text_eval",
                "value": {
                  "type": "determine",
                  "value": {
                    "$ARGUMENTS": null,
                    "$NAME": null
                  }
                }
              }
            }
          }
        }
      },
      {
        "op": "move",
        "from": "$.value['$NAME']",
        "path": "$.value['$FROM'].value.value['$VAL'].value.value['$NAME']",
        "reason": "move the raw name to the just inserted 'determine'."
      },
      {
        "op": "move",
        "from": "$.value['$ARGUMENTS']",
        "path": "$.value['$FROM'].value.value['$VAL'].value.value['$ARGUMENTS']",
        "reason": "move the arguments to the just inserted 'determine'."
      },
      {
        "op": "replace",
        "path": "$.type",
        "reason": "finally, rename #trying to #choose_value",
        "value": "choose_value"
      }
    ]
  }
]

function objField(obj, field) {
  return {
    "Object:dot:": [ 
      obj, {
        "At field:": field
      }
    ]
  }
}
// queries for title, score, etc.
const locationOfPlayer = {
  "FromText:": {
    "Determine:args:": [
      "location_of",
      {
        "Arg:from:": [
          "obj",
          {
            "FromText:": objField("story", "actor")
          }
        ]
      }
    ]
  }
};

export default {
  currentObjects: {
    "FromRecord:": {
      "Determine:args:": [
        "collect_objects",
        {
          "Arg:from:": [
            "obj", locationOfPlayer
          ]
        }
      ]
    }
  },
  currentScore: {
    "FromNum:": {
      "Num if:then:else:": [
        {
          "Is scene:": "scoring"
        }, 
        objField("story", "score"),
        -1
      ]
    }
  },
  currentTurn: {
    "FromNum:": {
      "Num if:then:else:": [
        {
          "Is scene:": "scoring"
        },
        objField("story", "turn count"),
        -1
      ]
    }
  },
  fabricate(text) {
    return {
      "FromExe:": {
        "Fabricate input:":text
      }
    }
  },
  locationName: {
    // fix: what do you mean this is insane?
    // it'd help if we could use the implicit pattern call decoder :/
    // maybe change "print name" to some "get name"
    // and, if possible, get rid of From(s)
    "FromText:": {
      "Buffer:": {
        "Determine:args:": [
          "print_name",
          {
            "Arg:from:": [
              "obj",  locationOfPlayer
            ]
          }
        ]
      }
    }
  },
   storyTitle: {
    "FromText:": objField("story", "title")
  }
}
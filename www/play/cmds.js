// queries for title, score, etc.
const locationOfPlayer = {
  "FromText:": {
    "Determine:args:": [
      "location_of", {
        "Arg:from:": [
          "obj", {
            "FromText:": {
              "Object:field:": ["story", "actor"]
            }
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
        "collect_objects", {
          "Arg:from:": [
            "obj", locationOfPlayer
          ]
        }
      ]
    }
  },
  currentScore: {
    "FromNumber:": {
      "Num if:then:else:": [
        { "Is domain:": "scoring" },
        {"Object:field:": ["story", "score"]},
        -1
      ]
    }
  },
  currentTurn:{
    "FromNumber:": {
      "Num if:then:else:": [
        { "Is domain:": "scoring" },
        {"Object:field:": ["story", "turn count"]},
        -1
      ]
    }
  },
  storyTitle: {
    "FromText:": {
      "Object:field:": ["story", "title"]
    }
  },
  locationName: {
    // fix: what do you mean this is insane?
    // it'd help if we could use the implicit pattern call decoder :/
    // maybe change "print name" to some "get name"
    // and, if possible, get rid of From(s)
      "FromText:": {
        "Buffers do:": {
          "Determine:args:": [
            "print_name", {
              "Arg:from:": [
                "obj",  locationOfPlayer
              ]
            }
          ]
        }
      }
   }
}

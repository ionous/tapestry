function localLang(make) {
  make.group("Story Statements", function() {
    make.flow("story", "{*paragraph}");
    make.flow("paragraph", "{*story_statement}", "Phrases");
    make.slot("story_statement", "Phrase");
    //
    make.flow("noun_statement", "story_statement", "{:lede} {*tail} {?summary}",
             "Noun statement: Describes people, places, or things.");

    make.flow("comment", "story_statement", "Note: {comment%lines}",
      "Comment: Information about the story not used by the story.")
  });

  make.group("Tests", function() {
    // "testing" is an interface, currently with once implementation type: TestOutput
    make.slot("testing", "Run a series of tests.");

    make.flow("test_statement", "story_statement",
      "Expect {test_name} to {expectation%test:testing}");

    make.flow("test_scene", "story_statement",
      "While testing {test_name}: {story}");

    make.flow("test_rule", "story_statement",
      "To test {test_name}: {do%hook:program_hook}");

    make.flow("test_output", "testing",
      "output {lines|quote}.",
      `Test Output: Expect that a test uses 'Say' to print some specific text.`);

    // would like just "<author's test name>" to be quoted, and not the current_test determiner.
    make.str("test_name", "{the test%current_test}, or {test name%test_name}");
  });

  make.group("Nouns", function() {
    make.flow("lede", "{nouns+named_noun|comma-and} {noun_phrase}.",
              "Leading statement: Describes one or more nouns.");

    make.flow("tail", "{pronoun} {noun_phrase}.",
             "Trailing statement: Adds details about the preceding noun or nouns.");

    // fix? change this into some sort of default pick of noun assignment.
    // make.flow("summary", "{The [summary] is:: %lines}");
    make.flow("summary", "The summary is: {summary%lines|quote}");

    make.opt("noun_phrase", "{kind_of_noun}, {noun_traits}, or {noun_relation}");

    // fix: think this should always be "are" never "is"
    make.flow("kind_of_noun", "{are_an} {trait*trait|comma-and} kind of {kind:singular_kind} {?noun_relation}");

    make.flow("noun_type",  "{an} {kind of%kinds:plural_kinds} noun");

    make.flow("named_noun", "object_eval", "{determiner} {name:noun_name}");

    make.str("determiner", "{a}, {an}, {the}, {our}, or {other determiner%determiner}",
      `Determiners: modify a word they are associated to designate specificity or, sometimes, a count.
        For instance: "some" fish hooks, "a" pineapple, "75" triangles, "our" Trevor.`  );

    make.str("noun_name",
      `Noun name: Some specific person, place, or thing; or, more rarely, a kind.
        Proper names are usually capitalized:  For example, maybe: 'Haruki', 'Jane', or 'Toronto'.
        Common names are usually not capitalized. For example, maybe: 'table', 'chair', or 'dog park'.
        A set of duplicate object uses their kind. For instance: twelve 'cats'.`);

    make.str("pronoun",  "{it}, {they}, or {pronoun}");

  });

  make.group("Patterns", function() {
    make.flow("pattern_decl", "story_statement",
       "The pattern {name:pattern_name|quote} determines {type:pattern_type}. {optvars?pattern_variables_tail} {about?comments}",
       `Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it.
  Each function in a given pattern has "guards" which determine whether the function applies in a particular situtation.`
     );

    make.flow("comments", "{lines|quote}");

    make.flow("pattern_variables_decl", "story_statement",
      "The pattern {pattern_name|quote} requires {+variable_decl|comma-and}.",
       `Declare pattern variables: Storage for values used during the execution of a pattern.`);

    make.flow("pattern_variables_tail", "It requires {+variable_decl|comma-and}.",
       `Pattern variables: Storage for values used during the execution of a pattern.`);

    make.opt("pattern_type", "an {activity:patterned_activity} or a {value:variable_type}");
    make.str("patterned_activity", "{an activity%activity}");
    make.str("pattern_name");

    make.flow("pattern_actions", "story_statement",
      "To {pattern name%name:pattern_name}: {?pattern_locals} {pattern_rules}",
      "Pattern Actions: Actions to take when using a pattern.");

    make.flow("pattern_rules", "{*pattern_rule}");
    make.flow("pattern_rule", `When {conditions are met%guard:bool_eval}, then: {do%hook:program_hook}`,
      "Rule");

    make.flow("pattern_locals", "{*local_decl}");
    make.flow("local_decl",  "Where {variable_decl} = {value%program_result}",
      "Local: local variables can use the parameters of a pattern to compute temporary values.");

    make.opt("program_hook", "flow an {activity} or return a {result:program_return}");

    // fix? activity and program_return both exist for the sake of appearance only.
    make.flow("program_return", "return {result:program_result}");

    make.opt("program_result", "a {simple value%primitive:primitive_func} or an {object:object_func}");
    make.opt("primitive_func", "{a number%number_eval}, {some text%text_eval}, {a true/false value%bool_eval}");
    make.flow("object_func", "an object named {name:text_eval}");
  });

  make.group("Relations", function() {
    make.flow("noun_relation",  "{?are_being} {relation} {nouns+named_noun|comma-and}");

    make.flow("relative_to_noun", "story_statement",
            "{relation} {nouns+named_noun} {are_being} {nouns+named_noun}.");

    make.str("relation");
  });

  make.group("Kinds", function() {
    make.flow("kinds_of_kind", "story_statement",
         "{plural_kinds} are a kind of {singular_kind}.");

    make.flow("kinds_possess_properties", "story_statement",
              "{plural_kinds} have {determiner} {property_phrase}.");

    make.str("singular_kind",
      `Kind: Describes a type of similar nouns.
For example: an animal, a container, etc.`);

    make.str("plural_kinds",
      `Kinds: The plural name of a type of similar nouns.
For example: animals, containers, etc.`);
  });

  make.group("Variables", function() {
    make.flow("variable_decl", "{type:variable_type} called {name:variable_name}");
    make.str("variable_name");

    make.opt("variable_type", "a {simple value%primitive:primitive_type}, an {object:object_type}, or {other value%ext:ext_type}");
    make.flow("object_type",  "{an} {kind of%kind:singular_kind}");
  });

  make.group("Traits", function() {
    make.flow("kinds_of_aspect", "story_statement", "{aspect} is a kind of value.");
    make.flow("aspect_traits", "story_statement", "{aspect} {trait_phrase}");
    make.flow("trait_phrase", "{are_either} {trait+trait|comma-or}.");

    make.flow("noun_traits", "{are_being} {trait+trait|comma-and}");
    make.str("aspect");
    make.str("trait");
  });

  make.group("Properties", function() {
    // ex. The description of the nets is xxx
    make.flow("noun_assignment", "story_statement",
            // "The {property} of {+noun} is the {[text]:: %lines}",
            "The {property} of {nouns+named_noun} is {the text%lines|summary}",
            "Noun Assignment: Assign text. Gives a noun one or more lines of text.");

    make.opt("property_phrase", "{primitive_phrase} or {aspect_phrase}");
    make.flow("aspect_phrase", "{aspect} {?optional_property}");
    make.flow("optional_property", "called {property}");

    make.flow("certainties", "story_statement",
              "{plural_kinds} {are_being} {certainty} {trait}.");

    make.str("are_either", "{can be%canbe} {are either%either}");

    make.str("certainty",  "{usually}, {always}, {seldom}, or {never}",
             "Certainty: Whether an trait applies to a kind of noun.");

    make.str("property");
  });

  // primitive types
  make.group("Values", function() {
    make.flow("primitive_phrase", "{primitive_type} called {property}");

    make.str("primitive_type", "{a number%number}, {some text%text}, or {a true/false value%bool}");
    make.opt("primitive_value", "{text%boxed_text} or {number%boxed_number}");

    // a list of numbers, a list of text, a record, or a list of records.
    make.opt("ext_type", "a list of {numbers:number_list}, a list of {text%text_list}, a {record:record_type} or a list of {records:record_list}.")

    make.flow("record_type",  "a record of {kind%kind:singular_kind}");
    make.flow("record_list",  "a list of {kind%kind:singular_kind} records");

    make.flow("boxed_text", "{text}");
    make.flow("boxed_number", "{number}");

    // constants
    make.str("text_list", "{a list of text%text_list}");
    make.str("number_list", "{a list of numbers%number_list}");

    make.str("bool", "{true} or {false}");
    make.str("text", "{text} or {empty}", `A sequence of characters of any length, all on one line.
Examples include letters, words, or short sentences.
Text is generally something displayed to the player.
See also: lines.`);

    // fix: bracket style links [] for see also?
    make.txt("lines", `A sequence of characters of any length spanning multiple lines.
Paragraphs are a prime example. Generally lines are some piece of the story that will be displayed to the player.
See also: text.`);
    make.num("number");
  });

 make.group("Helper Types", function() {
    make.str("an", "{a} or {an}");
    make.str("are_being",  "{are} or {is}");
    make.str("are_an",  "{are}, {are a%area}, {are an%arean}, {is}, {is a%isa}, {is an%isan}");
  });
}

function makeLang(make) {
  // read spec.js ( generated by iffy/cmd/spec/spec.go )
  make.group("Code", function() {
    spec.forEach((t)=> {
      make.newFromSpec(t);
    });
  });
  // read the local language
  make.group("Model", function() {
    localLang(make);
  });
}


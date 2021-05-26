function localLang(make) {
  // commands can only use other commands
  // so primitive tyoes need wrappers
  make.group("Value", function() {
    make.num("number");
    make.str("bool", "{true} or {false}");
    make.str("text", "{text} or {empty}", `A sequence of characters of any length, all on one line.
Examples include letters, words, or short sentences.
Text is generally something displayed to the player.
See also: lines.`);
    make.txt("lines", `A sequence of characters of any length spanning multiple lines.
Paragraphs are a prime example. Generally lines are some piece of the story that will be displayed to the player.
See also: text.`);

  //     "lines": {
  //   "uses": "flow",
  //   "spec": "here {_%lines:text}",
  //   "desc": "Specify one or more lines of text.",
  //   "group": "literals"
  // },

  });

  make.group("Model", function() {

    make.group("Statements", function() {
      make.flow("story", "{*paragraph}");
      make.flow("paragraph", "{*story_statement}", "Phrases");
      make.slot("story_statement", "Phrase");
      //
      make.flow("noun_statement", "story_statement", "{:lede} {*tail} {?summary}",
        "Declare a noun: Describes people, places, or things.");

      make.flow("comment", ["story_statement", "execute"], "Note: {comment%lines}",
        "Add a note: Information about the story for you and other authors.")
    });

    make.group("Tests", function() {
      // "testing" is an interface, currently with once implementation type: TestOutput
      make.slot("testing", "Run a series of tests.");

      make.flow("test_statement", "story_statement",
        "Expect {test_name} to {expectation%test:testing}",
        "Describe test results");

      make.flow("test_scene", "story_statement",
        "While testing {test_name}: {story}",
        "Create a scene for testing");

      make.flow("test_rule", "story_statement",
        "To test {test_name}: {do%hook:program_hook}",
        "Add actions to a test");

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

      make.swap("noun_phrase", "{is a kind%kind_of_noun}, {noun_traits}, or {noun_relation}");

      // fix: think this should always be "are" never "is"
      // fix: this shouldnt be "kind of", kind of declares a kind
      // ( but note singular vs. plural nouns phrases here along with are/is )
      // probably should have a switch for singular/ plural -- would be nice if are_an could look ahead and mutate with a custom filter maybe.
      make.flow("kind_of_noun", "{are_an} {*trait|comma-and} {kind:singular_kind} {?noun_relation}");

      make.flow("named_noun", "object_eval", "{determiner} {name:noun_name}");

      make.str("determiner", "{a}, {an}, {the}, {our}, or {other determiner%determiner}",
        `Determiners: modify a word they are associated to designate specificity or, sometimes, a count.
        For instance: "some" fish hooks, "a" pineapple, "75" triangles, "our" Trevor.`  );

      make.str("noun_name",
        `Noun name: Some specific person, place, or thing; or, more rarely, a kind.
        Proper names are usually capitalized:  For example, maybe: 'Haruki', 'Jane', or 'Toronto'.
        Common names are usually not capitalized. For example, maybe: 'table', 'chair', or 'dog park'.
        A set of duplicate object uses their kind. For instance: twelve 'cats'.`);

      make.str("pronoun", "{it}, {they}, or {pronoun}");

    });

    make.group("Events", function() {
      make.flow("action_decl", "story_statement", "Agents can {act%event:event_name} and {acting%action:action_name} requires {action_params}.",
        `Declare an activity: Activities help actors perform tasks: for instance, picking up or dropping items.
      Activities involve either the player or an npc and possibly one or two other objects.`);
      make.swap("action_params", "{one or more objects%common:common_action}, or {two similar objects%dual:paired_action}, or {nothing%none:abstract_action}")

      make.flow("common_action", "one {kind:singular_kind} ( the noun ) {?action_context}");
      const x = make.flow("action_context", "and one {kind:singular_kind} ( the other noun )");
      make.flow("paired_action", "two {kinds:plural_kinds} ( the noun and other noun )");

      make.str("action_name");

      make.flow("event_block", "story_statement", "For {the target%target:event_target} {handlers+event_handler}",
        `Declare event listeners: Listeners let objects in the game world react to changes before, during, or after they happen.`);
      make.flow("event_handler", "{event_phase} {the event%event:event_name} {with locals%locals?pattern_locals} do:{pattern_rules}");

      make.str("event_phase", "{before}, {during%while}, or {after}");
      make.str("event_name");
      make.swap("event_target", "the {kinds:plural_kinds} or {named_noun}");
    });

    make.group("Patterns", function() {
      make.flow("pattern_decl", "story_statement",
        "{type:pattern_type} determine {name:pattern_name|quote} {parameters%optvars?pattern_variables_tail} {?pattern_return} {about?comment}.",
        `Declare a pattern: A pattern is a bundle of functions which can either change the game world or provide information about it.
  Each function in a given pattern has "guards" which determine whether the function applies in a particular situtation.`
      );

      // fix: shouldnt this be called "parameters" -- maybe look at a few ebnfs to see what other languages use for naming things.
      make.flow("pattern_variables_decl", "story_statement",
        "The pattern {pattern_name|quote} requires {+variable_decl|comma-and}.",
        `Add parameters to a pattern: Values provided when calling a pattern.`);

      make.flow("pattern_variables_tail", "It requires {+variable_decl|comma-and}",
        `Pattern variables: Storage for values used during the execution of a pattern.`);

      make.str("pattern_type", "{patterns}, {actions}, {events}, or {another pattern type%pattern_type}");
      make.str("pattern_name");

      // fix: pattern return should be part of the declaration
      make.flow("pattern_actions", "story_statement",
        "To determine {pattern name%name:pattern_name} {?pattern_locals} {?pattern_return} do:{pattern_rules}",
        "Add actions to a pattern: Actions to take when using a pattern.");

      make.flow("pattern_rules", "{*pattern_rule}");
      make.flow("pattern_rule", `When {conditions are met%guard:bool_eval}{ continue%flags?pattern_flags}, then: {do%hook:program_hook}`,
        "Rule");

      make.str("pattern_flags", "{continue before%before}, {continue after%after}, {terminate}");

      make.flow("pattern_locals", "{+local_decl|comma-and}");

      make.flow("local_decl", " using {variable_decl} {starting as%value?local_init}",
        "Local: local variables can use the parameters of a pattern to compute temporary values.");

      make.flow("local_init", " starting as {value:assignment}",
        "Local: local variables can use the parameters of a pattern to compute temporary values.");

      make.swap("program_hook", "do {actions%activity}");

      make.flow("pattern_return", "returning {result:variable_decl}");
    });

    make.group("Relations", function() {
      make.flow("kind_of_relation", "story_statement", "{relation:relation_name} relates {relation_cardinality}");
      make.swap("relation_cardinality", "{one_to_one}, {one_to_many}, {many_to_one}, or {many_to_many}");
      make.flow("one_to_one", "one {kind:singular_kind} to one {other_kind:singular_kind}");
      make.flow("one_to_many", "one {kind:singular_kind} to many {kinds:plural_kinds}");
      make.flow("many_to_one", "many {kinds:plural_kinds} to one {kind:singular_kind}");
      make.flow("many_to_many", "many {kinds:plural_kinds} to many {other_kinds:plural_kinds}");

      make.flow("noun_relation", "{?are_being} {relation:relation_name} {nouns+named_noun|comma-and}");

      make.flow("relative_to_noun", "story_statement",
        "The {relation:relation_name} of {nouns+named_noun|comma-and} {are_being} {nouns+named_noun|comma-and}.",
        "Relate nouns to each other");

      // make.str("relation_name"); // also declared in rel.go
    });

    make.group("Kinds", function() {
      make.flow("kinds_of_kind", "story_statement",
        "{plural_kinds} are a kind of {singular_kind}.",
        "Declare a kind");

      make.flow("kinds_possess_properties", "story_statement",
        "{plural_kinds} have {+property_decl|comma-and}.",
        "Add properties to a kind");

      make.str("singular_kind",
        `Kind: Describes a type of similar nouns.
For example: an animal, a container, etc.`);

      make.str("plural_kinds",
        `Kinds: The plural name of a type of similar nouns.
For example: animals, containers, etc.`);
    });

    make.group("Records", function() {
      make.flow("kinds_of_record", "story_statement",
        "{records%record_plural} are a kind of record.",
        "Declare a record");

      make.flow("records_possess_properties", "story_statement",
        "{records%record_plural} have {+property_decl|comma-and}.",
        "Add properties to a record");

      make.str("record_singular",
        `Record: Describes a common set of properties.`);

      make.str("record_plural",
        `Records: The plural name for a record.`);
    });

    make.group("Traits", function() {
      make.flow("kinds_of_aspect", "story_statement", "{aspect} is a kind of value.",
        "Declare an aspect");
      make.flow("aspect_traits", "story_statement", "{aspect} {trait_phrase}",
        "Add traits to an aspect");
      make.flow("trait_phrase", "{are_either} {+trait|comma-or}.");

      make.flow("noun_traits", "{are_being} {+trait|comma-and}");
      make.str("aspect");
      make.str("trait");
    });

    make.group("Properties", function() {
      // ex. The description of the nets is xxx
      make.flow("noun_assignment", "story_statement",
        // "The {property} of {+noun} is the {[text]:: %lines}",
        "The {property} of {nouns+named_noun} is {the text%lines|quote}",
        "Assign text to a noun: Assign text. Gives a noun one or more lines of text.");

      make.flow("property_decl", "{an:determiner} {property} ( {property_type} {comment?lines} )");
      make.swap("property_type", "an {aspect%property_aspect}, {simple value%primitive:primitive_type}, or {other value%ext:ext_type}");

      make.str("property_aspect", "{an aspect%aspect}");

      make.flow("certainties", "story_statement",
        "{plural_kinds} {are_being} {certainty} {trait}.",
        "Give a kind a trait");

      make.str("are_either", "{can be%canbe} {are either%either}");

      make.str("certainty", "{usually}, {always}, {seldom}, or {never}",
        "Certainty: Whether an trait applies to a kind of noun.");

      make.str("property");
    });

    // primitive types

    make.group("Misc", function() {

      make.flow("record_type", "a record of {kind%kind:record_singular}");
      make.flow("record_list_type", "a list of {kind%kind:record_singular} records");

      make.flow("boxed_text", "{text}");
      make.flow("boxed_number", "{number}");

      make.str("ana", "{a} or {an}");
      make.str("are_being", "{are} or {is}");
      make.str("are_an", "{are}, {are a%area}, {are an%arean}, {is}, {is a%isa}, {is an%isan}");
      make.flow("variable_decl", "{an:determiner} {name:variable_name} ( {type:variable_type}  {comment?lines} )");
      // make.str("variable_name"); // also declared in core.

      make.swap("variable_type", "a {simple value%primitive:primitive_type}, an {object:object_type}, or {other value%ext:ext_type}");
      make.flow("object_type", "{an:ana} {kind of%kind:singular_kind}");

      make.str("primitive_type", "{a number%number}, {some text%text}, or {a true/false value%bool}");
      make.swap("primitive_value", "{text%boxed_text} or {number%boxed_number}");

      make.str("text_list_type", "{a list of text%list}");
      make.str("number_list_type", "{a list of numbers%list}");

      // a list of numbers, a list of text, a record, or a list of records.
      make.swap("ext_type", "a list of {numbers:number_list_type}, a list of {text%text_list_type}, a {record:record_type} or a list of {records:record_list_type}.")
    });
  });
}

function makeLang(make) {
  // read spec.js ( generated by iffy/cmd/spec/spec.go )
  make.group("Code", function() {
    spec.forEach((t) => {
      make.newFromSpec(t);
    });
  });
  // read the local language
  make.group("Model", function() {
    localLang(make);
  });
}

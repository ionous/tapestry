// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const Type_{{Pascal name}} = "{{name}}";

{{#if ../marshal}}
func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, ptr *{{Pascal name}}) (ret []byte,err error) {
  var b []byte
  if slat:= *ptr; slat != nil {
    b, err = slat.(jsonexp.DetailedMarshaler).MarshalDetailed(n)
  }
  if err == nil {
    ret, err= json.Marshal(jsonexp.Node{
      Type:  {{ScopeOf name}}Type_{{Pascal name}},
      Value: b,
    })
  }
  return
}

func {{Pascal name}}_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  var msg jsonexp.Node
  if e := json.Unmarshal(b, &msg); e != nil {
    err =  errutil.New(Type_{{Pascal name}}, "-", e)
  } else if contents:= msg.Value; len(contents) > 0 {
    var inner jsonexp.Node // peek to create the appropriate type
    if e := json.Unmarshal(contents, &inner); e != nil {
      err =  errutil.New("value of", Type_{{Pascal name}}, "-", e)
    } else if ptr, e := n.NewType(inner.Type); e != nil {
      err =  errutil.New(Type_{{Pascal name}}, "-", e)
    } else if imp, ok := ptr.(jsonexp.DetailedMarshaler); !ok {
      err =  errutil.New("casting slot", Type_{{Pascal name}}, "-", e)
    } else if e := imp.UnmarshalDetailed(n, contents); e != nil {
      err =  errutil.New("contents of", Type_{{Pascal name}}, "-", e)
    } else if fini,e := n.Finalize(ptr); e!= nil {
      err= e
    } else if store, ok:= fini.({{Pascal name}}); !ok {
      err= errutil.Fmt("couldnt store %T into %s", fini, Type_{{Pascal name}})
    } else {
      (*out) = store
    }
  }
  return
}
{{/if}}
{{/with}}
{{#if marshal}}
{{>repeat spec=this}}
{{/if}}
`;

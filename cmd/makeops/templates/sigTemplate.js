// sigTemplate.js
 module.exports = `

func (op* {{Pascal name}}) Marshal(n jsn.Marshaler) {
  {{Pascal name}}_Marshal(n, op)
}
`

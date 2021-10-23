// sigTemplate.js
 module.exports = `
func (op* {{Pascal name}}) Marshal(m jsn.Marshaler) error {
  return {{Pascal name}}_Marshal(m, op)
}
`

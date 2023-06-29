export function join(klass: string, joiner: string): string {
  const classes = [klass] as string[]
  if (joiner !== "") classes.push(joiner)
  return classes.join(" ")
}

export function joinIf(klass: string, condition: boolean, joiner: string): string {
  if (!condition) return klass
  return join(klass, joiner)
}

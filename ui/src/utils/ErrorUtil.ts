export function printError(err: unknown): Error {
  let error: Error
  if (err instanceof Error) {
    error = err
  } else {
    let errorString: string
    switch (typeof err) {
      case "string":
      case "object":
        errorString = JSON.stringify(err)
        break
      default:
        errorString = String(err)
    }
    error = new Error(errorString)
  }

  console.error(error) // eslint-disable-line no-console
  return error
}

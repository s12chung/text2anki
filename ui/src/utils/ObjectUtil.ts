// https://stackoverflow.com/a/56592365/1090482
export function pick<T extends object, K extends keyof T>(obj: T, ...keys: K[]) {
  return Object.fromEntries(keys.filter((key) => key in obj).map((key) => [key, obj[key]])) as Pick<
    T,
    K
  >
}

// https://stackoverflow.com/a/56592365/1090482
export function omit<T extends object, K extends keyof T>(obj: T, ...keys: K[]) {
  return Object.fromEntries(
    Object.entries(obj).filter(([key]) => !keys.includes(key as K)),
  ) as Omit<T, K>
}

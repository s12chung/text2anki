// https://stackoverflow.com/a/56592365/1090482
export const pick = <T extends object, K extends keyof T>(obj: T, ...keys: K[]) =>
  Object.fromEntries(keys.filter((key) => key in obj).map((key) => [key, obj[key]])) as Pick<T, K>

// https://stackoverflow.com/a/56592365/1090482
export const inclusivePick = <T extends object, K extends string | number | symbol>(
  obj: T,
  ...keys: K[]
) =>
  Object.fromEntries(keys.map((key) => [key, obj[key as unknown as keyof T]])) as {
    [key in K]: key extends keyof T ? T[key] : undefined
  }

// https://stackoverflow.com/a/56592365/1090482
export const omit = <T extends object, K extends keyof T>(obj: T, ...keys: K[]) =>
  Object.fromEntries(Object.entries(obj).filter(([key]) => !keys.includes(key as K))) as Omit<T, K>

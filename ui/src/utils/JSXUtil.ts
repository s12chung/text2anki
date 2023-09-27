import { DependencyList, ReactEventHandler, useCallback, useEffect, useState } from "react"

export function preventDefault(f: () => void): ReactEventHandler {
  return (e) => {
    e.preventDefault()
    f()
  }
}

export function useKeyDownEffect(
  keyDownHandler: (e: KeyboardEvent) => unknown,
  deps: DependencyList
) {
  /* eslint-disable react-hooks/exhaustive-deps */
  const wrappedHandler = useCallback(keyDownHandler, deps)
  useEffect(() => {
    window.addEventListener("keydown", wrappedHandler)
    return () => window.removeEventListener("keydown", wrappedHandler)
  }, [wrappedHandler])
}

export const useTimedState = (duration: number) => {
  const [value, setValue] = useState<boolean>(false)

  useEffect(() => {
    if (!value) {
      return () => {
        // do nothing
      }
    }
    const timer = setTimeout(() => setValue(false), duration)
    return () => clearTimeout(timer)
  }, [value, duration])
  return [value, setValue] as const
}

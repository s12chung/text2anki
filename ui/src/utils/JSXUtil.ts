import { DependencyList, ReactEventHandler, useCallback, useEffect } from "react"

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

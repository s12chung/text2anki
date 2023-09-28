import { newHash } from "./StringUtil.ts"
import { DependencyList, ReactEventHandler, useCallback, useEffect, useMemo, useState } from "react"

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

export function useTimedState(duration: number) {
  const [value, setValue] = useState<boolean>(false)

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    if (!value) return () => {}

    const timer = setTimeout(() => setValue(false), duration)
    return () => clearTimeout(timer)
  }, [value, duration])
  return [value, setValue] as const
}

export function useSafeSet(onSet: (val: boolean) => void, deps: DependencyList) {
  return useMemo(() => new SafeSet(onSet), deps)
}

export type Reset = () => void

export class SafeSet {
  private readonly resetMap: Record<string, Reset>
  constructor(private onSet: (val: boolean) => void) {
    this.resetMap = {}
  }

  addReset(reset: Reset): string {
    const key = newHash()
    this.resetMap[key] = reset
    return key
  }

  removeReset(key: string) {
    // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
    delete this.resetMap[key]
  }

  reset() {
    for (const key in this.resetMap) {
      if (!Object.hasOwn(this.resetMap, key)) continue
      this.resetMap[key]()
    }
  }

  safeSet(set: (val: boolean) => void, val: boolean) {
    this.reset()
    this.onSet(val)
    set(val)
  }
}

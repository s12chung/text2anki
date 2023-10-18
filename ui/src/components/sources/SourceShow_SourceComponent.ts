import React, { useMemo, useState } from "react"

interface StopKeyboard {
  stopKeyboardEvents: boolean
  setStopKeyboardEvents: (stop: boolean) => void
}
export const StopKeyboardContext = React.createContext<StopKeyboard>({
  stopKeyboardEvents: false,
  // eslint-disable-next-line no-empty-function
  setStopKeyboardEvents: () => {},
})

export function useStopKeyboard() {
  const [stopKeyboard, setStopKeyboard] = useState<boolean>(false)
  return useMemo<StopKeyboard>(
    () => ({ stopKeyboardEvents: stopKeyboard, setStopKeyboardEvents: setStopKeyboard }),
    [stopKeyboard],
  )
}

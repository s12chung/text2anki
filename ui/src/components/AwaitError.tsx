import { printError } from "../utils/ErrorUtil.ts"
import React from "react"
import { useAsyncError } from "react-router-dom"

const AwaitError: React.FC = () => {
  const error = printError(useAsyncError())
  return <div>Error: {error.message}</div>
}

export default AwaitError

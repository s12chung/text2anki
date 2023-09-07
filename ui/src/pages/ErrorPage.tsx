import { printError } from "../services/Format.ts"
import { useRouteError } from "react-router-dom"

const ErrorPage: React.FC = () => {
  const error = printError(useRouteError())
  return (
    <div className="m-std">
      <div className="text-xl">
        {error.name}: {error.message}
      </div>

      <div className="mt-std whitespace-pre-line">{error.stack}</div>
    </div>
  )
}

export default ErrorPage

import AwaitError from "./AwaitError.tsx"
import React from "react"
import { Await } from "react-router-dom"
import { AwaitResolveRenderFunction } from "react-router/dist/lib/components"

const AwaitWithFallback: React.FC<{
  readonly resolve: Promise<unknown>
  readonly children: AwaitResolveRenderFunction
}> = ({ resolve, children }) => (
  <React.Suspense fallback={<div>Loading...</div>}>
    <Await resolve={resolve} errorElement={<AwaitError />}>
      {children}
    </Await>
  </React.Suspense>
)

export default AwaitWithFallback

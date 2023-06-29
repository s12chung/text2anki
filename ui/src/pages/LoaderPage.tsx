import React from "react"
import { useLoaderData } from "react-router-dom"

export interface ILoaderPageProps<T> {
  Component: React.FC<{ data: T }>
}

const LoaderPage = <T extends object>({ Component }: ILoaderPageProps<T>) => {
  return <Component data={useLoaderData() as T} />
}

export default LoaderPage

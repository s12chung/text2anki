import "./index.css"
import HomePage from "./pages/HomePage"
import ApplicationLayout from "./pages/layouts/ApplicationLayout"
import sourceService from "./services/SourceService.ts"
import React from "react"
import ReactDOM from "react-dom/client"
import { createBrowserRouter, defer, RouterProvider } from "react-router-dom"

const router = createBrowserRouter([
  {
    path: "/",
    element: <ApplicationLayout />,
    children: [
      {
        path: "",
        element: <HomePage />,
        loader: () => defer({ sources: sourceService.list() }),
      },
    ],
  },
])

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)

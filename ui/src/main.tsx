import "./index.css"
import HomePage from "./pages/HomePage"
import SourceShowPage from "./pages/SourceShowPage.tsx"
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
      {
        path: "sources/:id",
        element: <SourceShowPage />,
        loader: ({ params }) => defer({ source: sourceService.get(params.id as string) }),
      },
    ],
  },
])

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
)

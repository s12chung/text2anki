import "./index.css"
import routes from "./routes.ts"
import React from "react"
import ReactDOM from "react-dom/client"
import { createBrowserRouter, RouterProvider } from "react-router-dom"

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <RouterProvider router={createBrowserRouter([routes])} />
  </React.StrictMode>
)

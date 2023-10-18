import Notifications from "./components/Notifications.tsx"
import ContextLayout from "./contexts/ContextLayout.tsx"
import "./index.css"
import routes from "./routes.ts"
import React from "react"
import ReactDOM from "react-dom/client"
import { createBrowserRouter, RouterProvider } from "react-router-dom"

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <ContextLayout>
      <RouterProvider router={createBrowserRouter([routes])} />
      <Notifications />
    </ContextLayout>
  </React.StrictMode>,
)

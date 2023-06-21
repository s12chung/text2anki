import * as SourceController from "./controllers/SourceController.ts"
import HomePage from "./pages/HomePage.tsx"
import SourceShowPage from "./pages/SourceShowPage.tsx"
import ApplicationLayout from "./pages/layouts/ApplicationLayout.tsx"
import { resources, route } from "./utils/RouterUtils.ts"
import { createElement } from "react"

const el = createElement

const routes = route("/", el(ApplicationLayout), {}, [
  route("", el(HomePage), { loader: SourceController.index }),

  resources("sources", SourceController, {
    show: el(SourceShowPage),
  }),
])

export default routes

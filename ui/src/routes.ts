import * as SourceController from "./controllers/SourcesController.ts"
import HomePage from "./pages/HomePage.tsx"
import SourceEditPage from "./pages/SourceEditPage.tsx"
import SourceNewPage from "./pages/SourceNewPage.tsx"
import SourceShowPage from "./pages/SourceShowPage.tsx"
import ApplicationLayout from "./pages/layouts/ApplicationLayout.tsx"
import FullLayout from "./pages/layouts/FullLayout.tsx"
import { IController, resources, route, withLayout } from "./utils/RouterUtils.ts"
import { createElement } from "react"

const el = createElement

const appLayoutSourceController = { get: SourceController.get } as IController

const routes = route("/", null, {}, [
  withLayout(el(ApplicationLayout), [
    route("", el(HomePage), { loader: SourceController.index }),

    resources("sources", appLayoutSourceController, {
      edit: el(SourceEditPage),
    }),
  ]),

  withLayout(el(FullLayout), [
    resources("sources", SourceController, {
      show: el(SourceShowPage),
      new: el(SourceNewPage),
    }),
  ]),
])

export default routes

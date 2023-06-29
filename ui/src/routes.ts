import * as SourceController from "./controllers/SourcesController.ts"
import HomePage from "./pages/HomePage.tsx"
import LoaderPage from "./pages/LoaderPage.tsx"
import SourceCreate from "./pages/components/source/SourceCreate.tsx"
import SourceEdit, { ISourceEditData } from "./pages/components/source/SourceEdit.tsx"
import SourceShow, { ISourceShowData } from "./pages/components/source/SourceShow.tsx"
import ApplicationLayout from "./pages/layouts/ApplicationLayout.tsx"
import FullLayout from "./pages/layouts/FullLayout.tsx"
import { pick } from "./utils/ObjectUtil.ts"
import { IController, resources, route, withLayout } from "./utils/RouterUtil.ts"
import { createElement } from "react"

const el = createElement

const appLayoutSourceController = pick(SourceController, "index", "create", "get") as IController

const routes = route("/", null, {}, [
  withLayout(el(ApplicationLayout), [
    route("", el(HomePage), { loader: SourceController.index }),

    resources("sources", appLayoutSourceController, {
      edit: el(LoaderPage<ISourceEditData>, { Component: SourceEdit }),
    }),
  ]),

  withLayout(el(FullLayout), [
    resources("sources", SourceController, {
      show: el(LoaderPage<ISourceShowData>, { Component: SourceShow }),
      new: el(SourceCreate),
    }),
  ]),
])

export default routes

import SourceCreate from "./components/source/SourceCreate.tsx"
import SourceEdit, { ISourceEditData } from "./components/source/SourceEdit.tsx"
import SourceShow, { ISourceShowData } from "./components/source/SourceShow.tsx"
import * as SourceController from "./controllers/SourcesController.ts"
import * as TermsController from "./controllers/TermsController.ts"
import HomePage from "./pages/HomePage.tsx"
import LoaderPage from "./pages/LoaderPage.tsx"
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
    route("terms/search", null, { loader: TermsController.search }),
  ]),

  withLayout(el(FullLayout), [
    resources("sources", SourceController, {
      show: el(LoaderPage<ISourceShowData>, { Component: SourceShow }),
      new: el(SourceCreate),
    }),
  ]),
])

export default routes

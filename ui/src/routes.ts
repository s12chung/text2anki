import PrePartListShow, {
  IPrePartListShowData,
} from "./components/pre_part_lists/PrePartListShow.tsx"
import SourceCreate from "./components/sources/SourceCreate.tsx"
import SourceEdit, { ISourceEditData } from "./components/sources/SourceEdit.tsx"
import SourceShow, { ISourceShowData } from "./components/sources/SourceShow.tsx"
import * as NotesController from "./controllers/NotesController.ts"
import * as PrePartListsController from "./controllers/PrePartListsController.ts"
import * as SourceController from "./controllers/SourcesController.ts"
import * as TermsController from "./controllers/TermsController.ts"
import HomePage from "./pages/HomePage.tsx"
import LoaderPage from "./pages/LoaderPage.tsx"
import ApplicationLayout from "./pages/layouts/ApplicationLayout.tsx"
import EmptyLayout from "./pages/layouts/EmptyLayout.tsx"
import FullLayout from "./pages/layouts/FullLayout.tsx"
import PrePartListDragAndDropLayout from "./pages/layouts/PrePartListDragAndDropLayout.tsx"
import { pick } from "./utils/ObjectUtil.ts"
import { IController, resources, route, withLayout } from "./utils/RouterUtil.ts"
import { createElement } from "react"

const el = createElement

const appLayoutSourceController: IController = pick(SourceController, "create", "edit")
const fullLayoutSourceController: IController = pick(SourceController, "get", "update", "destroy")

const routes = route("/", null, {}, [
  withLayout(el(PrePartListDragAndDropLayout), [
    route("", el(HomePage), { loader: SourceController.index }),
  ]),

  withLayout(el(ApplicationLayout), [
    resources("sources", appLayoutSourceController, {
      edit: el(LoaderPage<ISourceEditData>, { Component: SourceEdit }),
    }),
    route("terms/search", null, { loader: TermsController.search }),
    resources("notes", NotesController, {}),
  ]),

  withLayout(el(FullLayout), [
    resources("sources", fullLayoutSourceController, {
      show: el(LoaderPage<ISourceShowData>, { Component: SourceShow }),
      new: el(SourceCreate),
    }),
  ]),

  withLayout(el(EmptyLayout), [
    resources("sources", {}, {}, [
      resources("pre_part_lists", PrePartListsController, {
        show: el(LoaderPage<IPrePartListShowData>, { Component: PrePartListShow }),
      }),
    ]),
  ]),
])

export default routes

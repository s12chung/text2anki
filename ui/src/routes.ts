import NoteList, { INoteListData } from "./components/notes/NoteList.tsx"
import PrePartListShow, {
  IPrePartListShowData,
} from "./components/pre_part_lists/PrePartListShow.tsx"
import SourceCreate from "./components/sources/SourceCreate.tsx"
import SourceShow, { ISourceShowData } from "./components/sources/SourceShow.tsx"
import * as NotesController from "./controllers/NotesController.ts"
import * as PartsController from "./controllers/PartsController.ts"
import * as PrePartListsController from "./controllers/PrePartListsController.ts"
import * as SourceController from "./controllers/SourcesController.ts"
import * as TermsController from "./controllers/TermsController.ts"
import ErrorPage from "./pages/ErrorPage.tsx"
import HomePage from "./pages/HomePage.tsx"
import LoaderPage from "./pages/LoaderPage.tsx"
import ApplicationLayout from "./pages/layouts/ApplicationLayout.tsx"
import EmptyLayout from "./pages/layouts/EmptyLayout.tsx"
import FullLayout from "./pages/layouts/FullLayout.tsx"
import PrePartListDragAndDropLayout from "./pages/layouts/PrePartListDragAndDropLayout.tsx"
import { resources, route, withLayout } from "./utils/RouterUtil.ts"
import { createElement } from "react"

const el = createElement

const rootOptions = { errorElement: el(ErrorPage) }

const routes = route("/", null, rootOptions, [
  withLayout(el(PrePartListDragAndDropLayout), [
    route("", el(HomePage), { loader: SourceController.index }),
  ]),

  withLayout(el(ApplicationLayout), [
    route("terms/search", null, { loader: TermsController.search }),
    resources("notes", NotesController, {
      index: el(LoaderPage<INoteListData>, { Component: NoteList }),
    }),
  ]),

  withLayout(el(FullLayout), [
    resources("sources", SourceController, {
      show: el(LoaderPage<ISourceShowData>, { Component: SourceShow }),
      new: el(SourceCreate),
    }),
  ]),

  withLayout(el(EmptyLayout), [
    resources("sources", {}, {}, [
      resources("pre_part_lists", PrePartListsController, {
        show: el(LoaderPage<IPrePartListShowData>, { Component: PrePartListShow }),
      }),
      route(":sourceId", null, {}, [
        resources("parts", PartsController, {}, [
          route("multi", null, { action: PartsController.multi }),
        ]),
      ]),
    ]),
  ]),
])

export default routes

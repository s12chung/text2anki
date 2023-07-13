import { Status404, Status405 } from "./StatusUtil.ts"
import { ReactNode } from "react"
import { ActionFunction, ActionFunctionArgs, LoaderFunction, RouteObject } from "react-router-dom"

type IActionMap = Record<string, ActionFunction>

export function actionFunc(actionMap: IActionMap): ActionFunction {
  return async (args: ActionFunctionArgs) => {
    const fn = actionMap[args.request.method]
    if (typeof fn !== "function") {
      // eslint-disable-next-line @typescript-eslint/no-throw-literal
      throw new Response(`${args.request.method} not found`, Status405)
    }
    return fn(args)
  }
}

export function formData<T extends Record<keyof T, FormDataEntryValue>>(
  formData: FormData,
  ...keys: (keyof T)[]
): T {
  const obj = {} as T
  for (const key of keys) {
    obj[key] = formData.get(key as string) as T[keyof T]
  }
  return obj
}

export interface RouteOptions {
  loader?: LoaderFunction
  action?: ActionFunction
}

// eslint-disable-next-line max-params
export function route(
  path: string,
  element: ReactNode,
  options?: RouteOptions,
  children?: RouteObject[]
): RouteObject {
  const route = {
    element,
    children,
    ...options,
  } as RouteObject

  if (path === "") {
    route.index = true
  } else {
    route.path = path
  }

  return route
}

export function withLayout(element: ReactNode, children: RouteObject[]): RouteObject {
  return { element, children }
}

export interface IController {
  index?: LoaderFunction
  get?: LoaderFunction

  create?: ActionFunction
  update?: ActionFunction
  delete?: ActionFunction
}

export interface IElementMap {
  index?: ReactNode
  show?: ReactNode

  new?: ReactNode
  edit?: ReactNode
}

function resourceError(element: string, controllerMethod: string): Error {
  return new Error(`elements.${element} given, but no controller.${controllerMethod} exists`)
}

// eslint-disable-next-line max-params
export function resources(
  name: string,
  controller: IController,
  elements: IElementMap,
  children?: RouteObject[]
): RouteObject {
  const route = resourcesRoute(name, controller, elements)
  if (!route.children) route.children = []

  const idRoute = resourceRoute(controller, elements)
  if (idRoute) route.children.push(idRoute)

  const newRoute = newResourceRoute(elements)
  if (newRoute) route.children.push(newRoute)

  const editRoute = editResourceRoute(controller, elements)
  if (editRoute) route.children.push(editRoute)

  if (children) route.children.push(...children)

  return route
}

function path(url: string): string {
  return new URL(url).pathname.replace(/\/$/u, "")
}

function resourcesRoute(name: string, controller: IController, elements: IElementMap): RouteObject {
  const route = { path: name } as RouteObject

  if (elements.index) {
    if (!controller.index) throw resourceError("index", "index")
    route.element = elements.index
    route.loader = controller.index
  } else {
    route.loader = ({ request }) => {
      if (request.method === "GET" && path(request.url) === `/${name}`) {
        // eslint-disable-next-line @typescript-eslint/no-throw-literal
        throw new Response("path does not exist", Status404)
      }
      return null
    }
  }

  if (controller.create) {
    route.action = actionFunc({ POST: controller.create })
  }

  return route
}

function resourceRoute(controller: IController, elements: IElementMap): RouteObject | null {
  const route = { path: ":id" } as RouteObject
  if (elements.show) {
    if (!controller.get) throw resourceError("show", "get")
    route.element = elements.show
    route.loader = controller.get
  }

  const actionMap = {} as IActionMap
  if (controller.update) actionMap.PATCH = controller.update
  if (controller.delete) actionMap.DELETE = controller.delete

  if (Object.keys(actionMap).length !== 0) {
    route.action = actionFunc(actionMap)
  }

  return Object.keys(route).length === 1 ? null : route
}

function newResourceRoute(elements: IElementMap): RouteObject | null {
  if (!elements.new) {
    return null
  }
  return { path: "new", element: elements.new }
}

function editResourceRoute(controller: IController, elements: IElementMap): RouteObject | null {
  if (!elements.edit) {
    return null
  }
  if (!controller.get) {
    throw resourceError("edit", "get")
  }
  return { path: ":id/edit", element: elements.edit, loader: controller.get }
}

import { ReactNode } from "react"
import { ActionFunction, ActionFunctionArgs, LoaderFunction, RouteObject } from "react-router-dom"

type IActionMap = Record<string, ActionFunction>

export function actionFunc(actionMap: IActionMap): ActionFunction {
  return async (args: ActionFunctionArgs) => {
    const fn = actionMap[args.request.method]
    if (typeof fn !== "function") {
      // eslint-disable-next-line @typescript-eslint/no-throw-literal
      throw new Response(`${args.request.method} not found`, { status: 405 })
    }
    return fn(args)
  }
}

export function formData<T extends Record<keyof T, string>>(
  formData: FormData,
  ...keys: (keyof T)[]
): T {
  const obj = {} as T

  for (const key of keys) {
    const dataForKey = formData.get(key as string)
    if (typeof dataForKey === "string") obj[key] = dataForKey as T[keyof T]
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
  return {
    path,
    element,
    children,
    ...options,
  }
}

export interface IController {
  index?: LoaderFunction
  get?: LoaderFunction

  create?: ActionFunction
  update?: ActionFunction
}

export interface IElementMap {
  index?: ReactNode
  show?: ReactNode
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
  const resRoute = { path: name, children: [] } as RouteObject

  if (elements.index) {
    if (!controller.index) throw resourceError("index", "index")
    resRoute.element = elements.index
    resRoute.loader = controller.index
  }

  if (elements.show) {
    if (!controller.get) throw resourceError("show", "get")
    resRoute.children?.push(route(":id", elements.show, { loader: controller.get }))
  }
  if (elements.edit) {
    if (!controller.get) throw resourceError("edit", "get")
    resRoute.children?.push(route(":id/edit", elements.edit, { loader: controller.get }))
  }

  const actionMap = {} as IActionMap
  if (controller.create) {
    actionMap.POST = controller.create
  }
  if (controller.update) {
    actionMap.PATCH = controller.update
  }
  if (Object.keys(actionMap).length !== 0) {
    resRoute.action = actionFunc(actionMap)
  }

  if (children) resRoute.children?.push(...children)

  return resRoute
}

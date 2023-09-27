/* eslint-disable max-lines */
import { preventDefault } from "../utils/JSXUtil.ts"
import { Dialog, Transition } from "@headlessui/react"
import { XMarkIcon } from "@heroicons/react/24/outline"
import React, { Fragment } from "react"

export const SlideOverDialog: React.FC<{
  show: boolean
  onClose?: () => void
  leftNode?: React.ReactNode
  children: React.ReactNode
}> = ({ show, onClose, leftNode, children }) => {
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  if (!onClose) onClose = () => {}

  const dialogWidth = "42rem"

  return (
    <Transition.Root show={show} as={Fragment}>
      {/* eslint-disable-next-line react/forbid-component-props */}
      <Dialog as="div" className="relative z-10" onClose={onClose}>
        {leftNode ? (
          <div
            style={{
              width: `calc(100vw - ${dialogWidth})`,
            }}
          >
            {leftNode}
          </div>
        ) : (
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" />
        )}

        <div className={`fixed inset-0 overflow-hidden${leftNode ? " w-0" : ""}`}>
          <div className="absolute inset-0 overflow-hidden">
            <div className="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10 sm:pl-16">
              <Transition.Child
                as={Fragment}
                enter="transform transition ease-in-out duration-500 sm:duration-700"
                enterFrom="translate-x-full"
                enterTo="translate-x-0"
                leave="transform transition ease-in-out duration-500 sm:duration-700"
                leaveFrom="translate-x-0"
                leaveTo="translate-x-full"
              >
                <Dialog.Panel
                  /* eslint-disable-next-line react/forbid-component-props */
                  className="pointer-events-auto"
                  /* eslint-disable-next-line react/forbid-component-props */
                  style={{
                    width: dialogWidth,
                  }}
                >
                  {/* eslint-disable-next-line react/jsx-max-depth */}
                  <div className="h-full overflow-y-scroll bg-white shadow-xl">{children}</div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </div>
      </Dialog>
    </Transition.Root>
  )
}

SlideOverDialog.defaultProps = {
  // eslint-disable-next-line no-undefined
  leftNode: undefined,
  // eslint-disable-next-line no-undefined
  onClose: undefined,
}

export const SlideOverHeader: React.FC<{
  title: string
  subtitle?: string
  onClose?: () => void
}> = ({ title, subtitle, onClose }) => {
  return (
    <div className="bg-gray-50 px-4 py-6 sm:px-6">
      <div className="flex items-start justify-between space-x-3">
        <div className="space-y-1">
          {/* eslint-disable-next-line react/forbid-component-props */}
          <Dialog.Title className="text-base font-semibold leading-6 text-gray-900">
            {title}
          </Dialog.Title>
          {Boolean(subtitle) && <p className="text-sm text-gray-500">{subtitle}</p>}
        </div>
        {onClose ? (
          <div className="flex h-7 items-center">
            <a href="#" className="a-btn" onClick={preventDefault(onClose)}>
              <span className="sr-only">Close panel</span>
              <XMarkIcon className="h-6 w-6" aria-hidden="true" />
            </a>
          </div>
        ) : null}
      </div>
    </div>
  )
}

SlideOverHeader.defaultProps = {
  // eslint-disable-next-line no-undefined
  onClose: undefined,
  subtitle: "",
}

export default {
  Dialog: SlideOverDialog,
  Header: SlideOverHeader,
}

import NotificationsContext, {
  Notifier,
  Notification,
  NotificationType,
} from "../contexts/NotificationsContext.ts"
import { preventDefault } from "../utils/JSXUtil.ts"
import { Transition } from "@headlessui/react"
import { InformationCircleIcon, XCircleIcon, XMarkIcon } from "@heroicons/react/20/solid"
import { CheckCircleIcon } from "@heroicons/react/24/outline"
import React, { Fragment, MouseEventHandler, useContext, useEffect, useState } from "react"

const Notifications: React.FC = () => {
  const { notifications, setNotifications } = useContext<Notifier>(NotificationsContext)
  const [notification, setNotification] = useState<Notification | null>(null)

  useEffect(() => {
    if (!notification) return
    const timer = setTimeout(() => setNotification(null), 7000)
    // eslint-disable-next-line consistent-return
    return () => clearTimeout(timer)
  }, [notification])

  useEffect(() => {
    if (notifications.length === 0) return
    setNotification(notifications[0])
    setNotifications(notifications.slice(1))
  }, [notifications, setNotifications])

  return notification ? (
    <NotificationWrapper
      key={notification.createdAt}
      notification={notification}
      onClose={preventDefault(() => setNotification(null))}
    />
  ) : null
}

const NotificationWrapper: React.FC<{
  readonly notification: Notification
  readonly onClose: MouseEventHandler<HTMLAnchorElement>
}> = ({ notification, onClose }) => {
  return (
    <div
      aria-live="assertive"
      className="pointer-events-none fixed inset-0 flex items-end justify-end px-4 py-6 sm:p-6"
    >
      <div className="flex w-full flex-col items-center space-y-4 sm:items-end">
        <Transition
          show
          as={Fragment}
          enter="transform ease-out duration-300 transition"
          enterFrom="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
          enterTo="translate-y-0 opacity-100 sm:translate-x-0"
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="pointer-events-auto w-full max-w-sm overflow-hidden rounded-lg bg-white shadow-lg ring-1 ring-black ring-opacity-5">
            <NotificationContents notification={notification} onClose={onClose} />
          </div>
        </Transition>
      </div>
    </div>
  )
}

const NotificationContents: React.FC<{
  readonly notification: Notification
  readonly onClose: MouseEventHandler<HTMLAnchorElement>
}> = ({ notification, onClose }) => {
  return (
    <div className="p-4 flex items-start">
      <div className="flex-shrink-0">
        {notification.type === NotificationType.Info && (
          <InformationCircleIcon className="h-6 w-6" aria-hidden="true" />
        )}
        {notification.type === NotificationType.Success && (
          <CheckCircleIcon className="h-6 w-6 text-green-400" aria-hidden="true" />
        )}
        {notification.type === NotificationType.Error && (
          <XCircleIcon className="h-6 w-6 text-red-600" aria-hidden="true" />
        )}
      </div>
      <div className="ml-3 w-0 flex-1 pt-0.5">
        <p className="text-sm font-medium text-gray-900">{notification.message}</p>
      </div>
      <div className="ml-4 flex flex-shrink-0">
        <a href="#" className="inline-flex a-btn" onClick={onClose}>
          <XMarkIcon className="h-5 w-5" aria-hidden="true" />
        </a>
      </div>
    </div>
  )
}

export default Notifications

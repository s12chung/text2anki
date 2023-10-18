import NotificationsContext, { Notification, Notifier } from "./NotificationsContext.ts"
import React, { useMemo, useState } from "react"

const ContextLayout: React.FC<{ readonly children: React.ReactNode }> = ({ children }) => {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const notifier = useMemo<Notifier>(
    () => new Notifier(notifications, setNotifications),
    [notifications],
  )

  return <NotificationsContext.Provider value={notifier}>{children}</NotificationsContext.Provider>
}

export default ContextLayout

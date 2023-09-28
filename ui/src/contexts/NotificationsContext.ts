import React, { Dispatch, SetStateAction } from "react"

export enum NotificationType {
  Info,
  Success,
  Error,
}

export interface Notification {
  message: string
  type: NotificationType
  createdAt: number
}

export class Notifier {
  constructor(
    public notifications: Notification[],
    public setNotifications: Dispatch<SetStateAction<Notification[]>>
  ) {}

  private notify(message: string, type: NotificationType) {
    const notification = {
      message,
      type,
      createdAt: Date.now(),
    }
    this.setNotifications([...this.notifications, notification])
  }

  info = (message: string) => {
    this.notify(message, NotificationType.Info)
  }

  success = (message: string) => {
    this.notify(message, NotificationType.Success)
  }

  error = (message: string) => {
    this.notify(message, NotificationType.Error)
  }
}

// eslint-disable-next-line @typescript-eslint/no-empty-function
const NotificationsContext = React.createContext<Notifier>(new Notifier([], () => {}))
export default NotificationsContext

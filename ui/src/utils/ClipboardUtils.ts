import { printError } from "./ErrorUtil.ts"

export function imageToClipboard(imageUrl: string) {
  const image = new Image()
  image.crossOrigin = "anonymous"
  const canvas = document.createElement("canvas")
  const ctx = canvas.getContext("2d")
  if (!ctx) throw new Error("error getting canvas context")

  image.onload = () => {
    canvas.width = image.width
    canvas.height = image.height
    ctx.drawImage(image, 0, 0)
    canvas.toBlob((blob) => {
      if (!blob) {
        printError("failed to load canvas image for clipboard")
        return
      }
      navigator.clipboard
        .write([new ClipboardItem({ [blob.type]: blob })])
        .catch((err) => printError(err))
    }, "image/png")
  }
  image.src = imageUrl
}

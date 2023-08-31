export function imageToClipboard(imageUrl: string): Promise<void> {
  const image = new Image()
  image.src = imageUrl
  image.crossOrigin = "anonymous"
  const canvas = document.createElement("canvas")
  const ctx = canvas.getContext("2d")

  return new Promise<void>((resolve, reject) => {
    if (!ctx) {
      reject(new Error("error getting canvas context"))
      return
    }

    image.onload = () => {
      canvas.width = image.width
      canvas.height = image.height
      ctx.drawImage(image, 0, 0)
      canvas.toBlob((blob) => {
        if (!blob) {
          reject(new Error("failed to load canvas image for clipboard"))
          return
        }
        navigator.clipboard
          .write([new ClipboardItem({ [blob.type]: blob })])
          .then(resolve)
          .catch(reject)
      }, "image/png")
    }
  })
}

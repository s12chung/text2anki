export interface PrePart {
  imageUrl: string
  audioUrl: string
}
const PrePartEmpty = Object.freeze<PrePart>({
  imageUrl: "",
  audioUrl: "",
})

export interface PrePartList {
  id: string
  preParts: PrePart[]
}
export const PrePartListEmpty = Object.freeze<PrePartList>({
  id: "",
  preParts: [PrePartEmpty],
})

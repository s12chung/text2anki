import ApplicationService from "./ApplicationService.ts"
import { Http, requestInit } from "./Format.ts"
import { Source, SourceEmpty } from "./models/Source.ts"

export interface PartData {
  text: string
  translation: string
}
export const PartDataEmpty = Object.freeze<PartData>({
  text: "",
  translation: "",
})

export interface PartCreateMultiData {
  prePartListId: string
  parts: PartData[]
}
export const PartCreateMultiDataEmpty = Object.freeze<PartCreateMultiData>({
  prePartListId: "",
  parts: [PartDataEmpty],
})

class PartsService extends ApplicationService {
  protected pathPrefix = "/sources"

  async create(sourceId: number, data: PartData): Promise<Source> {
    return this.fetch(`/${sourceId}/parts`, SourceEmpty, requestInit(Http.POST, data))
  }

  async multi(sourceId: number, data: PartCreateMultiData): Promise<Source> {
    return this.fetch(`/${sourceId}/parts/multi`, SourceEmpty, requestInit(Http.POST, data))
  }

  async update(
    sourceId: number | string,
    partIndex: number | string,
    data: PartData
  ): Promise<Source> {
    return this.fetch(`/${sourceId}/parts/${partIndex}`, SourceEmpty, requestInit(Http.PATCH, data))
  }

  async destroy(sourceId: number | string, partIndex: number | string): Promise<Source> {
    return this.fetch(`/${sourceId}/parts/${partIndex}`, SourceEmpty, requestInit(Http.DELETE))
  }
}

export const partsService = new PartsService()

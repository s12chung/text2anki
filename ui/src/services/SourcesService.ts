import ApplicationService from "./ApplicationService"
import { Http, requestInit } from "./Format.ts"
import { PartCreateMultiData, PartCreateMultiDataEmpty } from "./PartsService.ts"
import { Source, SourceEmpty } from "./models/Source.ts"

export interface CreateSourceData extends PartCreateMultiData {
  name: string
  reference: string
}
export const CreateSourceDataEmpty = Object.freeze<CreateSourceData>({
  name: "",
  reference: "",
  ...PartCreateMultiDataEmpty,
})

export interface UpdateSourceData {
  name: string
  reference: string
}
export const UpdateSourceDataEmpty = Object.freeze<UpdateSourceData>({
  name: "",
  reference: "",
})

class SourcesService extends ApplicationService {
  protected pathPrefix = "/sources"

  async index(): Promise<Source[]> {
    return this.fetch("", [SourceEmpty])
  }

  async get(id: string | number): Promise<Source> {
    return this.fetch(`/${id}`, SourceEmpty)
  }

  async create(data: CreateSourceData): Promise<Source> {
    return this.fetch("", SourceEmpty, requestInit(Http.POST, data))
  }

  async update(id: string | number, data: UpdateSourceData): Promise<Source> {
    return this.fetch(`/${id}`, SourceEmpty, requestInit(Http.PATCH, data))
  }

  async destroy(id: string | number): Promise<Source> {
    return this.fetch(`/${id}`, SourceEmpty, requestInit(Http.DELETE))
  }
}

export const sourcesService = new SourcesService()

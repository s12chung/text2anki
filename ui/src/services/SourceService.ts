import ApplicationService from "./ApplicationService"

export interface Source {
  id: number
  name: string
}

class SourceService extends ApplicationService {
  protected pathPrefix = "/sources"

  async list(): Promise<Source[]> {
    return (await this.fetch()) as Source[]
  }
}

const sourceService = new SourceService()

export default sourceService

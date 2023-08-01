import { PrePartList } from "../../services/PrePartListsService.ts"
import AwaitError from "../AwaitError.tsx"
import React from "react"
import { Await } from "react-router-dom"

export interface IPrePartListSourceCreateData {
  prePartList: Promise<PrePartList>
}
interface IPrePartListSourceCreateProps {
  data: IPrePartListSourceCreateData
}

const PrePartListSourceCreate: React.FC<IPrePartListSourceCreateProps> = ({ data }) => {
  return (
    <React.Suspense fallback={<div>Loading....</div>}>
      <Await resolve={data.prePartList} errorElement={<AwaitError />}>
        {(prePartList: PrePartList) => prePartList.id}
      </Await>
    </React.Suspense>
  )
}

export default PrePartListSourceCreate

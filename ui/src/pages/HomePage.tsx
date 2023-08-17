import SourceCreateMini from "../components/sources/SourceCreateMini.tsx"
import SourceList, { ISourceListData } from "../components/sources/SourceList"
import { useLoaderData } from "react-router-dom"

const HomePage: React.FC = () => {
  const data = useLoaderData() as ISourceListData
  return (
    <>
      <SourceCreateMini />
      <SourceList data={data} />
    </>
  )
}

export default HomePage

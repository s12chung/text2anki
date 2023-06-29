import SourceCreateMini from "./components/source/SourceCreateMini.tsx"
import SourceList, { ISourceListData } from "./components/source/SourceList"
import { useLoaderData } from "react-router-dom"

const HomePage: React.FC = () => {
  const data = useLoaderData() as ISourceListData
  return (
    <div>
      <SourceCreateMini />
      <SourceList data={data} />
    </div>
  )
}

export default HomePage

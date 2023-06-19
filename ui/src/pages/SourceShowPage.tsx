import SourceShow, {
  ISourceShowData,
} from "../components/source/SourceShow.tsx"
import { useLoaderData } from "react-router-dom"

const SourceShowPage: React.FC = () => {
  const data = useLoaderData() as ISourceShowData
  return <SourceShow data={data} />
}

export default SourceShowPage

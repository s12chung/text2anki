import { Link, Outlet } from "react-router-dom"

const FullLayout: React.FC = () => {
  return (
    <div className="m-std">
      <div className="my-std">
        <Link to="/">text2anki</Link>
      </div>
      <div>
        <Outlet />
      </div>
    </div>
  )
}

export default FullLayout

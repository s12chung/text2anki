import { Link, Outlet } from "react-router-dom"

const ApplicationLayout: React.FC = () => {
  return (
    <div className="grid-std">
      <div className="py-std">
        <Link to="/">text2anki</Link>
      </div>
      <div>
        <Outlet />
      </div>
    </div>
  )
}

export default ApplicationLayout

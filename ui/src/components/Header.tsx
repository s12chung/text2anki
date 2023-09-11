import React from "react"
import { Link } from "react-router-dom"

const Header: React.FC = () => {
  return (
    <div className="py-std flex-std">
      <div className="flex-grow">
        <Link to="/">text2anki</Link>
      </div>
      <Link to="/notes">Notes</Link>
    </div>
  )
}

export default Header

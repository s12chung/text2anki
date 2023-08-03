import React from "react"
import { Link } from "react-router-dom"

const Header: React.FC = () => {
  return (
    <div className="py-std">
      <Link to="/">text2anki</Link>
    </div>
  )
}

export default Header

import { Form } from "react-router-dom"

const SourceCreate: React.FC = () => {
  return (
    <Form action="/sources" method="post">
      <textarea name="text" placeholder="You may also drag and drop here." />
      <button type="submit">Submit</button>
    </Form>
  )
}

export default SourceCreate
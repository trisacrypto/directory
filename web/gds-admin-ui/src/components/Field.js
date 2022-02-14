import React from 'react'
import { Form } from "react-bootstrap"

const Input = React.forwardRef(({ register, name, ...rest }, ref) => <Form.Control type="text" ref={ref} {...register(name)} {...rest} />)
const Select = ({ register, name, ...rest }) => <Form.Select {...register(name)} {...rest} />
const Switch = ({ register, name, ...rest }) => <Form.Check {...register(name)} {...rest} />

const Field = () => {
    return <></>
}

Field.Input = Input
Field.Select = Select
Field.Switch = Switch

export default Field
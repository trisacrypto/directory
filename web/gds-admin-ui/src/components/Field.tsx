import React from 'react';
import { Form } from 'react-bootstrap';
import { UseFormRegister } from 'react-hook-form';

const BsControl = Form.Control;
const BsSelect = Form.Select;
const BsCheck = Form.Check;

type BaseProps = {
    register: UseFormRegister<any>;
    name: string;
};

type InputProps = BaseProps & React.ComponentPropsWithoutRef<typeof BsControl>;
type SelectProps = BaseProps & React.ComponentPropsWithoutRef<typeof BsSelect>;
type SwitchProps = BaseProps & React.ComponentPropsWithoutRef<typeof BsCheck>;

const Input = React.forwardRef(({ register, name, ...rest }: InputProps, ref) => (
    <Form.Control type="text" {...register(name)} {...rest} />
));
const Select = ({ register, name, ...rest }: SelectProps) => <Form.Select {...register(name)} {...rest} />;
const Switch = ({ register, name, ...rest }: SwitchProps) => <Form.Check {...register(name)} {...rest} />;

const Field = () => <></>;

Field.Input = Input;
Field.Select = Select;
Field.Switch = Switch;

export default Field;

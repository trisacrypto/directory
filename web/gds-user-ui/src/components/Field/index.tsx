import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import React from 'react';
import { RegisterOptions } from 'react-hook-form';

type InputProps = {
  register: RegisterOptions;
  name: string;
};

const Input = React.forwardRef<any, any>(({ register, name, ...rest }, ref) => (
  <InputFormControl type="text" ref={ref} {...rest} />
));

Input.displayName = 'Input';

const Select: React.FC<any> = ({ register, name, ...rest }) => (
  <SelectFormControl controlId={name} {...rest} />
);

const Field = () => {
  return <></>;
};

Field.Input = Input;
Field.Select = Select;

export default Field;

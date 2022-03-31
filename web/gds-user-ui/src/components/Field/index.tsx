import InputFormControl from 'components/ui/InputFormControl';
import SelectFormControl from 'components/ui/SelectFormControl';
import React from 'react';
import { RegisterOptions } from 'react-hook-form';

type InputProps = {
  register: RegisterOptions;
  name: string;
};

const Input = React.forwardRef<any, any>(({ register, name, ...rest }, ref) => (
  <InputFormControl type="text" ref={ref} {...rest} {...register(name)} />
));

Input.displayName = 'Input';

type SelectProps = {
  register: RegisterOptions;
  name: string;
};

const Select = React.forwardRef<SelectProps, any>(({ register, name, ...rest }, ref) => (
  <SelectFormControl controlId={name} ref={ref} {...rest} {...register(name)} />
));

Select.displayName = 'Select';

const Field = () => {
  return <></>;
};

Field.Input = Input;
Field.Select = Select;

export default Field;

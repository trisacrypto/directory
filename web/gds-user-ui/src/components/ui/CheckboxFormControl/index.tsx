import {
  FormControl,
  FormHelperText,
  FormLabel,
  Checkbox,
  CheckboxGroup,
  CheckboxProps,
  FormErrorMessage
} from '@chakra-ui/react';
import React from 'react';
interface _FormControlProps extends CheckboxProps {
  formHelperText?: string | React.ReactNode;
  controlId: string;
  name: string;
  isDisabled?: boolean;
  isRequired?: boolean;
}

const CheckboxFormControl = React.forwardRef<any, _FormControlProps>(
  ({ formHelperText, controlId, name, isDisabled, isRequired, isInvalid, ...rest }, ref) => {
    return (
      <FormControl isInvalid={isInvalid}>
        <FormLabel htmlFor={controlId}>{name}</FormLabel>
        <CheckboxGroup>
          <Checkbox
            id={controlId}
            name={name}
            isDisabled={isDisabled}
            isRequired={isRequired}
            ref={ref}
            {...rest}
          />
        </CheckboxGroup>
        {!isInvalid ? (
          <FormHelperText>{formHelperText}</FormHelperText>
        ) : (
          <FormErrorMessage>{formHelperText}</FormErrorMessage>
        )}
      </FormControl>
    );
  }
);

CheckboxFormControl.displayName = 'CheckboxFormControl';
export default CheckboxFormControl;

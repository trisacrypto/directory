import {
  FormControl,
  FormControlProps,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Switch,
  SwitchProps
} from '@chakra-ui/react';
import React from 'react';

interface _FormControlProps extends FormControlProps {
  formHelperText?: string;
  controlId: string;
  label?: string;
  inputProps?: SwitchProps;
  name?: string;
  error?: string;
}

const SwitchFormControl = React.forwardRef<any, _FormControlProps>(
  ({ inputProps, controlId, label, formHelperText, isInvalid, name, error, ...rest }, ref) => {
    return (
      <FormControl display="flex" alignItems="center" gap={2} isInvalid={!!error} {...rest}>
        <Switch id={controlId} {...inputProps} name={name} ref={ref} />
        <FormLabel htmlFor={controlId} mb={0}>
          {label}
        </FormLabel>
        {!isInvalid ? (
          <FormHelperText position="absolute" top={4}>
            {formHelperText}
          </FormHelperText>
        ) : (
          <FormErrorMessage position="absolute" top={4}>
            {formHelperText}
          </FormErrorMessage>
        )}
      </FormControl>
    );
  }
);

SwitchFormControl.displayName = 'SwitchFormControl';

export default SwitchFormControl;

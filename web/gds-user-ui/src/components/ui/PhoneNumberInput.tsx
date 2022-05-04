import 'react-phone-number-input/style.css';
import PhoneInput, { Props } from 'react-phone-number-input';
import {
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Input,
  InputProps,
  useColorModeValue
} from '@chakra-ui/react';
import React from 'react';

interface _Props extends Props<InputProps> {
  formHelperText?: string;
  controlId: string;
  label?: string;
  onChange: (arg: any) => void;
}

const PhoneNumberInput = React.forwardRef<any, _Props>(
  ({ formHelperText, isInvalid, controlId, label, onChange, ...props }, ref) => {
    const inputColorMode = useColorModeValue('#E3EBEF', undefined);

    return (
      <FormControl isInvalid={isInvalid}>
        <FormLabel htmlFor={controlId}>{label}</FormLabel>
        <PhoneInput
          ref={ref}
          onChange={onChange}
          background={inputColorMode}
          inputComponent={Input}
          borderRadius={0}
          limitMaxLength
          {...props}
        />
        {!isInvalid ? (
          <FormHelperText>{formHelperText}</FormHelperText>
        ) : (
          <FormErrorMessage>{formHelperText}</FormErrorMessage>
        )}
      </FormControl>
    );
  }
);

PhoneNumberInput.displayName = 'PhoneNumberInput';

export default PhoneNumberInput;

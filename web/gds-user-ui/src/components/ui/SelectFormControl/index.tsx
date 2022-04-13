import {
  FormHelperText,
  FormLabel,
  FormControl,
  useColorModeValue,
  FormErrorMessage
} from '@chakra-ui/react';
import { ChakraStylesConfig, GroupBase, OptionsOrGroups, Select, Props } from 'chakra-react-select';
import React from 'react';

interface _FormControlProps extends Props {
  formHelperText?: string;
  controlId: string;
  label?: string;
  name?: string;
  placeholder?: string;
  isDisabled?: boolean;
  options?: OptionsOrGroups<unknown, GroupBase<unknown>>;
}

const SelectFormControl = React.forwardRef<any, _FormControlProps>(
  (
    {
      label,
      formHelperText,
      controlId,
      placeholder,
      name,
      isDisabled,
      options,
      isMulti,
      isInvalid,
      ...rest
    },
    ref
  ) => {
    const bgColorMode = useColorModeValue('#E3EBEF', undefined);
    const chakraStyles: ChakraStylesConfig = {
      control: (provided, state) => ({
        ...provided,
        background: bgColorMode,
        borderRadius: 0
      })
    };

    return (
      <FormControl isInvalid={isInvalid}>
        <FormLabel htmlFor={controlId}>{label}</FormLabel>
        <Select
          name={name}
          id={controlId}
          placeholder={placeholder}
          chakraStyles={chakraStyles}
          options={options}
          isDisabled={isDisabled}
          isMulti={isMulti as any}
          {...rest}
          ref={ref}
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

SelectFormControl.displayName = 'SelectFormControl';

export default SelectFormControl;

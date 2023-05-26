import {
  FormHelperText,
  FormLabel,
  FormControl,
  useColorModeValue,
  FormErrorMessage
} from '@chakra-ui/react';
import { ChakraStylesConfig, GroupBase, OptionsOrGroups, Select, Props } from 'chakra-react-select';
import React, { ReactNode } from 'react';

interface _FormControlProps extends Props {
  formHelperText?: any;
  controlId: string;
  label?: string | ReactNode;
  name?: string;
  placeholder?: string;
  isDisabled?: boolean;
  defaultValue?: any;
  [key: string]: any;
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
      defaultValue,
      ...rest
    },
    ref
  ) => {
    const bgColorMode = useColorModeValue('#E3EBEF', undefined);
    const chakraStyles: ChakraStylesConfig = {
      control: (provided) => ({
        ...provided,
        background: bgColorMode,
        borderRadius: 0
      }),
      option: (provided, state) => ({
        ...provided,
        color: state.isSelected ? 'gray.500' : undefined
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
          defaultValue={defaultValue}
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

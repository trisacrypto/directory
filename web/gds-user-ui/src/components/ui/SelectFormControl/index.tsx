import {
  FormHelperText,
  FormLabel,
  FormControl,
  useColorModeValue,
  FormErrorMessage,
} from "@chakra-ui/react";
import {
  ChakraStylesConfig,
  GroupBase,
  OptionsOrGroups,
  Select,
  Props,
} from "chakra-react-select";

interface _FormControlProps extends Props {
  formHelperText?: string;
  controlId: string;
  label?: string;
  name?: string;
  placeholder?: string;
  options?: OptionsOrGroups<unknown, GroupBase<unknown>>;
}

const SelectFormControl: React.FC<_FormControlProps> = ({
  label,
  formHelperText,
  controlId,
  placeholder,
  name,
  options,
  isMulti,
  isInvalid,
  ...rest
}) => {
  const bgColorMode = useColorModeValue("#E3EBEF", undefined);
  const chakraStyles: ChakraStylesConfig = {
    control: (provided, state) => ({
      ...provided,
      background: bgColorMode,
      borderRadius: 0,
    }),
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
        isMulti={isMulti as any}
        {...rest}
      />
      {!isInvalid ? (
        <FormHelperText>{formHelperText}</FormHelperText>
      ) : (
        <FormErrorMessage>{formHelperText}</FormErrorMessage>
      )}
    </FormControl>
  );
};

export default SelectFormControl;

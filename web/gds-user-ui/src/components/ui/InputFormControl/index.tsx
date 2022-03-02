import {
  FormControl as CkFormControl,
  FormControlProps,
  FormHelperText,
  FormLabel,
  Input,
  InputProps,
  useColorModeValue,
  FormErrorMessage,
} from "@chakra-ui/react";

interface _FormControlProps extends FormControlProps {
  formHelperText?: string;
  controlId: string;
  label?: string;
  inputProps?: InputProps;
  name?: string;
  error?: string;
}

const InputFormControl: React.FC<_FormControlProps> = ({
  label,
  formHelperText,
  controlId,
  inputProps,
  name,
  isInvalid,
}) => {
  const inputColorMode = useColorModeValue("#E3EBEF", undefined);

  return (
    <CkFormControl isInvalid={isInvalid}>
      <FormLabel htmlFor={controlId}>{label}</FormLabel>
      <Input
        name={name}
        id={controlId}
        background={inputColorMode}
        borderRadius={0}
        {...inputProps}
      />
      {!isInvalid ? (
        <FormHelperText>{formHelperText}</FormHelperText>
      ) : (
        <FormErrorMessage>{formHelperText}</FormErrorMessage>
      )}
    </CkFormControl>
  );
};

export default InputFormControl;

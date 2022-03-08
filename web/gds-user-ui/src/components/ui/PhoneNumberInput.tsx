import "react-phone-number-input/style.css";
import PhoneInput, { Props } from "react-phone-number-input";
import {
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Input,
  InputProps,
  useColorModeValue,
} from "@chakra-ui/react";

interface _Props extends Props<InputProps> {
  formHelperText?: string;
  controlId: string;
  label?: string;
}

const PhoneNumberInput: React.FC<_Props> = ({
  formHelperText,
  isInvalid,
  controlId,
  label,
  ...props
}) => {
  const inputColorMode = useColorModeValue("#E3EBEF", undefined);

  return (
    <FormControl isInvalid={isInvalid}>
      <FormLabel htmlFor={controlId}>{label}</FormLabel>
      <PhoneInput
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
};

export default PhoneNumberInput;

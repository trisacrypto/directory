import {
  FormControl,
  FormControlProps,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Switch,
  SwitchProps,
} from "@chakra-ui/react";

interface _FormControlProps extends FormControlProps {
  formHelperText?: string;
  controlId: string;
  label?: string;
  inputProps?: SwitchProps;
  name?: string;
  error?: string;
}

const SwitchFormControl: React.FC<_FormControlProps> = ({
  inputProps,
  controlId,
  label,
  formHelperText,
  isInvalid,
  name,
  error,
  ...rest
}) => {
  return (
    <FormControl
      display="flex"
      alignItems="center"
      gap={2}
      isInvalid={!!error}
      {...rest}
    >
      <Switch id={controlId} {...inputProps} name={name} />
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
};

export default SwitchFormControl;

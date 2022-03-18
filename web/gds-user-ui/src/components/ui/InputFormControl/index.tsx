import {
  FormControl as CkFormControl,
  FormControlProps,
  FormHelperText,
  FormLabel,
  Input,
  InputProps,
  useColorModeValue,
  FormErrorMessage,
  InputRightElement,
  Button,
  InputGroup
} from '@chakra-ui/react';

interface _FormControlProps extends FormControlProps {
  formHelperText?: string;
  controlId: string;
  label?: string;
  inputProps?: InputProps;
  name?: string;
  error?: string;
  type?: React.HTMLInputTypeAttribute;
  hasBtn?: boolean;
  value?: string;
  setBtnName?: string;
  handleUserUpdate?: () => void;
}

const InputFormControl: React.FC<_FormControlProps> = ({
  label,
  formHelperText,
  controlId,
  inputProps,
  name,
  isInvalid,
  type = 'text',
  hasBtn,
  value,
  setBtnName,
  handleUserUpdate
}) => {
  const inputColorMode = useColorModeValue('#E3EBEF', undefined);

  return (
    <CkFormControl isInvalid={isInvalid}>
      <FormLabel htmlFor={controlId}>{label}</FormLabel>
      <InputGroup>
        <Input
          name={name}
          id={controlId}
          background={inputColorMode}
          borderRadius={0}
          type={type}
          value={value}
          {...inputProps}
        />
        {hasBtn && (
          <InputRightElement width="4.5rem">
            <Button
              h="1.75rem"
              bg={'transparent'}
              color={'blue'}
              size="sm"
              onClick={handleUserUpdate}>
              {setBtnName || 'Change'}
            </Button>
          </InputRightElement>
        )}
      </InputGroup>
      {!isInvalid ? (
        <FormHelperText>{formHelperText}</FormHelperText>
      ) : (
        <FormErrorMessage>{formHelperText}</FormErrorMessage>
      )}
    </CkFormControl>
  );
};

export default InputFormControl;

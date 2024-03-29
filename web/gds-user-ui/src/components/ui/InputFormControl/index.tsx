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
import React from 'react';

export interface _FormControlProps extends Omit<FormControlProps, 'label'> {
  formHelperText?: string | React.ReactNode;
  controlId: string;
  label?: React.ReactNode;
  inputProps?: InputProps;
  shouldResetValue?: boolean;
  name?: string;
  error?: string;
  type?: React.HTMLInputTypeAttribute;
  hasBtn?: boolean;
  value?: string;
  setBtnName?: string;
  isRequired?: boolean;
  inputRef?: React.RefObject<HTMLInputElement>;
  onValueChange?: any;
  handleFn?: () => void;
  [key: string]: any;
  isHidden?: boolean;
  isRequiredField?: boolean;
}

const InputFormControl = React.forwardRef<any, _FormControlProps>(
  (
    {
      label,
      formHelperText,
      controlId,
      inputProps,
      name,
      isInvalid,
      type = 'text',
      hasBtn,
      inputRef,
      setBtnName,
      handleFn,
      onChange,
      isDisabled,
      isRequired,
      placeholder,
      isRequiredField,
      ...rest
    },
    ref
  ) => {
    const inputColorMode = useColorModeValue('#E3EBEF', undefined);

    const handleMouseScroll = (e: React.WheelEvent<HTMLInputElement>) => {
      // Disable Mouse scrolling
      if (e.currentTarget.type === 'number') {
        e.currentTarget.blur();
      }
    };

    return (
      <CkFormControl isInvalid={isInvalid}>
        <FormLabel htmlFor={controlId}>
          {label}
          {isRequiredField && <span style={{ color: 'red' }}>*</span>}
        </FormLabel>
        <InputGroup>
          <Input
            name={name}
            id={controlId}
            background={inputColorMode}
            borderRadius={0}
            type={type}
            ref={inputRef || ref}
            onChange={onChange}
            onWheel={handleMouseScroll}
            isDisabled={isDisabled}
            isRequired={isRequired}
            placeholder={placeholder}
            {...inputProps}
            {...rest}
          />
          {hasBtn && (
            <InputRightElement width="4.5rem" height={'100%'}>
              <Button size="sm" onClick={handleFn}>
                {setBtnName || 'Change'}
              </Button>
            </InputRightElement>
          )}
        </InputGroup>
        {!isInvalid ? (
          <FormHelperText>{formHelperText}</FormHelperText>
        ) : (
          <FormErrorMessage role="alert" data-testid="error-message">
            {formHelperText}
          </FormErrorMessage>
        )}
      </CkFormControl>
    );
  }
);

InputFormControl.displayName = 'InputFormControl';

export default InputFormControl;

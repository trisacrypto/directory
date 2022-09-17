// create a custom toast component form chackra ui
import React from 'react';
import { useToast, UseToastOptions } from '@chakra-ui/react';

interface IToastProps {
  description?: string;
  status?: string;
  isClosable?: boolean;
  position?: string;
  duration: number;
}
const Toast = ({
  description,
  position,
  status,
  isClosable,
  duration
}: IToastProps & UseToastOptions) => {
  const toast = useToast();
  return toast({
    status: status ? status : 'error',
    description,
    isClosable: isClosable ? isClosable : true,
    position: position ? position : 'top-right',
    duration: duration ? duration : 5000
  });
};

export default Toast;

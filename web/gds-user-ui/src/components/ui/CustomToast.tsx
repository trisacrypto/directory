// create a custom toast component form chackra ui

import { useToast, UseToastOptions } from '@chakra-ui/react';

interface IToastProps {
  title: string;
  description?: string;
  status?: string;
  isClosable?: boolean;
  position?: string;
  duration?: number;
}
const CustomToast = ({
  description,
  position,
  status,
  isClosable,
  duration,
  title
}: IToastProps & UseToastOptions) => {
  const toast = useToast();
  toast({
    title,
    status: status ? status : 'error',
    description,
    isClosable: isClosable ? isClosable : true,
    position: position ? position : 'top-right',
    duration: duration ? duration : 9000
  });
  return null;
};

export default CustomToast;

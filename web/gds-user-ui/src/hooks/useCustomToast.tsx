// custom toast hook from chakra ui toast

import { useToast, UseToastOptions } from '@chakra-ui/react';

interface IToastProps {
  description?: string;
  status?: UseToastOptions['status'];
  isClosable?: boolean;
  position?: UseToastOptions['position'];
  duration?: number;
}

const useCustomToast = () => {
  const toast = useToast();
  const Toast = (props: IToastProps) => {
    return toast({
      status: props.status ? props.status : 'error',
      description: props.description,
      isClosable: props.isClosable ? props.isClosable : true,
      position: props.position ? props.position : 'top-right',
      duration: props.duration ? props.duration : 5000
    });
  };
  return Toast;
};

export default useCustomToast;

// this is custom hook for toast with chakra ui and react

import {
  useToast as useChakraToast,
  UseToastOptions,
  ToastPositionWithLogical
} from '@chakra-ui/react';

interface ToastProps {
  id?: string;
  title: string;
  description?: string;
  status?: 'info' | 'warning' | 'success' | 'error';
  duration?: number;
  isClosable?: boolean;
  position?: ToastPositionWithLogical;
}

interface ToastOption {
  id: string;
}

const useToast = ({ id }: ToastOption) => {
  const { title, description, status, duration, isClosable, position } = {
    id: id || 'toast',
    position: 'top'
  } as ToastProps;
  const toast = useChakraToast();

  const toastProps: UseToastOptions = {
    id,
    title,
    description,
    status,
    duration,
    isClosable,
    position
  };

  return toast(toastProps);
};

export default useToast;

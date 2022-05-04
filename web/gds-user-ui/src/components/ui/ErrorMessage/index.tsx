import React, { useState, useEffect } from 'react';
import { Box, Alert, AlertIcon, AlertDescription, CloseButton } from '@chakra-ui/react';
type TError = {
  message: string;
  handleClose?: () => void;
};

const ErrorMessage: React.FC<TError> = ({ message, handleClose }) => {
  return (
    <Box my={4}>
      <Alert status="error" borderRadius={4}>
        <AlertIcon />
        <AlertDescription>{message}</AlertDescription>
        <CloseButton position="absolute" right="8px" top="8px" onClick={handleClose} />
      </Alert>
    </Box>
  );
};
export default ErrorMessage;
